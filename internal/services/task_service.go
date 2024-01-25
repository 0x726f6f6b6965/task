package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0x726f6f6b6965/task/internal/helper"
	"github.com/0x726f6f6b6965/task/internal/utils"
	pbTask "github.com/0x726f6f6b6965/task/protos/task/v1"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	TaskID  string = "taskID"
	SortSet string = "sortSet"
)

type taskService struct {
	pbTask.UnimplementedTaskServiceServer
	sequencer   utils.Generator
	redisClient *redis.Client
	logger      *zap.Logger
}

// CreateTask - create a task
func (service *taskService) CreateTask(ctx context.Context, req *pbTask.CreateTaskRequest) (*pbTask.Task, error) {
	if helper.IsEmpty(req.GetName()) {
		return nil, helper.RequiredFieldErr("name is empty", "name")
	}

	if _, ok := pbTask.Status_name[int32(req.Status)]; !ok {
		return nil, helper.InvalidErr("status invalid", "status", req.Status)
	}

	var (
		id string
	)

	task := &pbTask.Task{
		Name:   req.GetName(),
		Status: req.Status,
	}
	seq, _ := service.sequencer.Next()
	id = seq.String()

	// check id exist
	// usually the id won't repeat
	exist := service.redisClient.Exists(ctx, fmt.Sprintf("%s:%s", TaskID, id)).Val()

	if exist != 0 || helper.IsEmpty(id) {
		service.logger.Error("CreateTask attempt to create id error", zap.Any("request", req))
		return nil, helper.InternalErr("please try again later")
	}

	task.Id = id

	data, err := json.Marshal(task)
	if err != nil {
		service.logger.Error("CreateTask unmarshal error", zap.Error(err))
		return nil, helper.InternalErr("unmarshal error")
	}

	err = service.redisClient.Set(ctx, fmt.Sprintf("%s:%s", TaskID, id), data, -1).Err()
	if err != nil {
		service.logger.Error("CreateTask redis set error", zap.Error(err))
		return nil, helper.InternalErr("redis set error")
	}

	err = service.redisClient.ZAdd(ctx, SortSet, redis.Z{
		Score:  seq.Float64(),
		Member: id,
	}).Err()

	if err != nil {
		service.logger.Error("CreateTask redis zadd error", zap.Error(err))
		return nil, helper.InternalErr("redis zadd error")
	}
	return task, nil
}

// DeleteTask - delete a task by id
func (service *taskService) DeleteTask(ctx context.Context, req *pbTask.DeleteTaskRequest) (*emptypb.Empty, error) {
	if helper.IsEmpty(req.GetId()) {
		return nil, helper.RequiredFieldErr("id is empty", "id")
	}
	exist := service.redisClient.Exists(ctx, fmt.Sprintf("%s:%s", TaskID, req.GetId())).Val()
	if exist == 0 {
		return nil, helper.NotFoundErr("task not found", "id", req.GetId())
	}
	ok, err := service.redisClient.Del(ctx, fmt.Sprintf("%s:%s", TaskID, req.GetId())).Result()
	if err != nil {
		service.logger.Error("DeleteTask redis del error", zap.String("id", req.GetId()), zap.Error(err))
		return nil, helper.InternalErr("redis del error")
	}

	if ok != 1 {
		service.logger.Error("DeleteTask del fail", zap.String("id", req.GetId()))
		return nil, helper.InternalErr("please try again later")
	}
	ok, err = service.redisClient.ZRem(ctx, SortSet, req.GetId()).Result()
	if err != nil {
		service.logger.Error("DeleteTask redis zrem error", zap.String("id", req.GetId()), zap.Error(err))
		return nil, helper.InternalErr("redis zrem error")
	}

	if ok != 1 {
		service.logger.Error("DeleteTask zrem fail", zap.String("id", req.GetId()))
		return nil, helper.InternalErr("please try again later")
	}
	return &emptypb.Empty{}, nil
}

// GetTask - get task information
func (service *taskService) GetTask(ctx context.Context, req *pbTask.GetTaskRequest) (*pbTask.Task, error) {
	if helper.IsEmpty(req.GetId()) {
		return nil, helper.RequiredFieldErr("id is empty", "id")
	}

	result, err := service.redisClient.Get(ctx, fmt.Sprintf("%s:%s", TaskID, req.GetId())).Bytes()
	if err != nil {
		if errors.Is(redis.Nil, err) {
			return nil, helper.NotFoundErr("task not found", "id", req.GetId())
		}
		service.logger.Error("GetTask redis get error", zap.Error(err))
		return nil, helper.InternalErr("redis get error")
	}
	resp := &pbTask.Task{}
	err = json.Unmarshal(result, resp)
	if err != nil {
		service.logger.Error("GetTask unmarshal error", zap.Error(err))
		return nil, helper.InternalErr("unmarshal error")
	}
	return resp, nil
}

// GetTaskList - get a list of task information
func (service *taskService) GetTaskList(ctx context.Context, req *pbTask.GetTaskListRequest) (*pbTask.GetTaskListResponse, error) {
	var (
		start string = "-"
		size  int64  = 25
		token        = utils.NewPageToken("", 0)
	)

	if !helper.IsEmpty(req.PageToken) {
		token, err := utils.GetPageTokenByString(req.PageToken)
		// if we get the wrong token, ignore it.
		if err == nil {
			start = "(" + token.GetID()
			size = token.GetSize()
		}
	}

	if req.PageSize != 0 {
		size = int64(req.PageSize)
	}
	keys, err := service.redisClient.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:    SortSet,
		ByLex:  true,
		Start:  start,
		Stop:   "+",
		Offset: 0,
		Count:  size,
	}).Result()

	if err != nil {
		service.logger.Error("GetTaskList redis zrange error", zap.Error(err))
		return nil, helper.InternalErr("redis zrange error")
	}
	resp := &pbTask.GetTaskListResponse{
		Tasks: []*pbTask.Task{},
	}
	for _, key := range keys {
		bytes, _ := service.redisClient.Get(ctx, fmt.Sprintf("%s:%s", TaskID, key)).Bytes()
		task := &pbTask.Task{}
		err = json.Unmarshal(bytes, task)
		if err != nil {
			service.logger.Error("GetTaskList unmarshal error",
				zap.String("key", fmt.Sprintf("%s:%s", TaskID, key)), zap.Error(err))
			continue
		}
		resp.Tasks = append(resp.Tasks, task)
	}

	if len(keys) >= int(size) {
		token.SetID(keys[len(keys)-1])
		token.SetSize(size)
		resp.NextToken = token.GetToken()
	}

	return resp, nil
}

// UpdateTask - update a task information by id
func (service *taskService) UpdateTask(ctx context.Context, req *pbTask.UpdateTaskRequest) (*pbTask.Task, error) {
	if helper.IsEmpty(req.GetId()) {
		return nil, helper.RequiredFieldErr("id is empty", "id")
	}
	if req.Task == nil {
		return nil, helper.RequiredFieldErr("task is empty", "task")
	}
	if _, ok := pbTask.Status_name[int32(req.Task.Status)]; !ok {
		return nil, helper.InvalidErr("status invalid", "status", req.Task.Status)
	}
	data, err := service.redisClient.Get(ctx, fmt.Sprintf("%s:%s", TaskID, req.GetId())).Bytes()
	if err != nil {
		if errors.Is(redis.Nil, err) {
			return nil, helper.NotFoundErr("task not found", "id", req.GetId())
		}
		service.logger.Error("UpdateTask redis get error", zap.Error(err))
		return nil, helper.InternalErr("redis get error")
	}
	task := &pbTask.Task{}
	err = json.Unmarshal(data, task)
	if err != nil {
		service.logger.Error("UpdateTask unmarshal error", zap.Error(err))
		return nil, helper.InternalErr("unmarshal error")
	}

	for _, key := range req.UpdateMask.GetPaths() {
		switch key {
		case "task.name":
			task.Name = req.Task.Name
		case "task.status":
			task.Status = req.Task.Status
		}
	}

	data, err = json.Marshal(task)
	if err != nil {
		service.logger.Error("UpdateTask unmarshal error", zap.Error(err))
		return nil, helper.InternalErr("unmarshal error")
	}

	err = service.redisClient.Set(ctx, fmt.Sprintf("%s:%s", TaskID, req.Id), data, -1).Err()
	if err != nil {
		service.logger.Error("UpdateTask redis set error", zap.Error(err))
		return nil, helper.InternalErr("redis set error")
	}
	return task, nil
}

func NewTaskService(generator utils.Generator, redisClient *redis.Client, logger *zap.Logger) pbTask.TaskServiceServer {
	return &taskService{
		redisClient: redisClient,
		sequencer:   generator,
		logger:      logger,
	}
}
