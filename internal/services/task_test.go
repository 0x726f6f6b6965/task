package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/0x726f6f6b6965/task/internal/helper"
	"github.com/0x726f6f6b6965/task/internal/utils"
	pbTask "github.com/0x726f6f6b6965/task/protos/task/v1"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	rClient *redis.Client
	rmock   redismock.ClientMock
	service pbTask.TaskServiceServer
	mockG   utils.Generator
	num     *big.Int
	ctx     context.Context
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	rClient, rmock = redismock.NewClientMock()
	logger, _ := zap.NewDevelopment()
	mockG = &mockGenerator{}
	num = big.NewInt(time.Now().UnixMilli())
	service = NewTaskService(mockG, rClient, logger)
	ctx = context.Background()
	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	rClient.Close()
	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestCreateTask(t *testing.T) {
	var (
		req  = &pbTask.CreateTaskRequest{Name: "test-name", Status: 1}
		g, _ = mockG.Next()
		task = &pbTask.Task{
			Id:     g.String(),
			Name:   req.Name,
			Status: req.Status,
		}
		data, _ = json.Marshal(task)
		key     = fmt.Sprintf("%s:%s", TaskID, g.String())
	)

	rmock.ExpectExists(key).SetVal(0)
	rmock.ExpectEval(helper.AddTask, []string{key, SortSet}, data, g.String()).RedisNil()

	resp, err := service.CreateTask(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, task, resp)
}

func TestCreateTaskEmptyName(t *testing.T) {
	var (
		req = &pbTask.CreateTaskRequest{Status: 1}
	)

	_, err := service.CreateTask(context.Background(), req)
	assert.Contains(t, err.Error(), "name is empty")
}

func TestCreateTaskIdExist(t *testing.T) {
	var (
		req  = &pbTask.CreateTaskRequest{Name: "test-name", Status: 1}
		g, _ = mockG.Next()
		key  = fmt.Sprintf("%s:%s", TaskID, g.String())
	)

	rmock.ExpectExists(key).SetVal(1)

	_, err := service.CreateTask(context.Background(), req)
	assert.Contains(t, err.Error(), "please try again later")
}

func TestCreateTaskInvalidStatus(t *testing.T) {
	var (
		req = &pbTask.CreateTaskRequest{Name: "test-name", Status: 3}
	)

	_, err := service.CreateTask(context.Background(), req)
	assert.Contains(t, err.Error(), "status invalid")
}

func TestGetTask(t *testing.T) {
	var (
		g, _   = mockG.Next()
		except = &pbTask.Task{
			Id:     g.String(),
			Name:   "test-name",
			Status: 1}
		data, _ = json.Marshal(except)
		key     = fmt.Sprintf("%s:%s", TaskID, g.String())
	)
	rmock.ExpectGet(key).SetVal(string(data))

	resp, err := service.GetTask(context.Background(), &pbTask.GetTaskRequest{Id: g.String()})
	assert.Nil(t, err)
	assert.Equal(t, except, resp)
}

func TestGetTaskEmptyId(t *testing.T) {
	_, err := service.GetTask(context.Background(), &pbTask.GetTaskRequest{})
	assert.Contains(t, err.Error(), "id is empty")
}

func TestGetTaskNotFound(t *testing.T) {
	var (
		g, _ = mockG.Next()
		key  = fmt.Sprintf("%s:%s", TaskID, g.String())
	)
	rmock.ExpectGet(key).RedisNil()

	_, err := service.GetTask(context.Background(), &pbTask.GetTaskRequest{Id: g.String()})
	assert.Contains(t, err.Error(), "task not found")
}

func TestDeleteTask(t *testing.T) {
	var (
		g, _ = mockG.Next()
		key  = fmt.Sprintf("%s:%s", TaskID, g.String())
	)
	rmock.ExpectExists(key).SetVal(1)
	rmock.ExpectEval(helper.DeleteTask, []string{key, SortSet}, g.String()).RedisNil()

	_, err := service.DeleteTask(context.Background(), &pbTask.DeleteTaskRequest{Id: g.String()})
	assert.Nil(t, err)
}

func TestDeleteTaskEmptyId(t *testing.T) {
	_, err := service.DeleteTask(context.Background(), &pbTask.DeleteTaskRequest{})
	assert.Contains(t, err.Error(), "id is empty")
}

func TestDeleteTaskNotFound(t *testing.T) {
	var (
		g, _ = mockG.Next()
		key  = fmt.Sprintf("%s:%s", TaskID, g.String())
	)
	rmock.ExpectExists(key).SetVal(0)

	_, err := service.DeleteTask(context.Background(), &pbTask.DeleteTaskRequest{Id: g.String()})
	assert.Contains(t, err.Error(), "task not found")
}

func TestUpdateTask(t *testing.T) {
	var (
		g, _ = mockG.Next()
		task = &pbTask.Task{
			Id:     g.String(),
			Name:   "test-name",
			Status: 1}
		data, _ = json.Marshal(task)
		key     = fmt.Sprintf("%s:%s", TaskID, g.String())
		req     = &pbTask.UpdateTaskRequest{
			Id:         g.String(),
			Task:       task,
			UpdateMask: &fieldmaskpb.FieldMask{},
		}
	)
	req.Task.Name = "update-test-name"
	req.UpdateMask.Paths = append(req.UpdateMask.Paths, "task.name")
	updateData, _ := json.Marshal(req.Task)
	rmock.ExpectGet(key).SetVal(string(data))
	rmock.ExpectSet(key, updateData, -1).SetVal("OK")
	resp, err := service.UpdateTask(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, req.Task, resp)
}

func TestUpdateTaskEmptyId(t *testing.T) {
	_, err := service.UpdateTask(context.Background(), &pbTask.UpdateTaskRequest{})
	assert.Contains(t, err.Error(), "id is empty")
}

func TestUpdateTaskEmptyTask(t *testing.T) {
	g, _ := mockG.Next()
	_, err := service.UpdateTask(context.Background(), &pbTask.UpdateTaskRequest{Id: g.String()})
	assert.Contains(t, err.Error(), "task is empty")
}

func TestUpdateTaskNotFound(t *testing.T) {
	var (
		g, _ = mockG.Next()
		task = &pbTask.Task{
			Id:     g.String(),
			Name:   "test-name",
			Status: 1}
		req = &pbTask.UpdateTaskRequest{
			Id:         g.String(),
			Task:       task,
			UpdateMask: &fieldmaskpb.FieldMask{},
		}
		key = fmt.Sprintf("%s:%s", TaskID, g.String())
	)

	rmock.ExpectGet(key).RedisNil()

	_, err := service.UpdateTask(context.Background(), req)
	assert.Contains(t, err.Error(), "task not found")
}

func TestUpdateTaskInvalidStatus(t *testing.T) {
	var (
		req = &pbTask.UpdateTaskRequest{Id: "test-id", Task: &pbTask.Task{Status: 3}}
	)

	_, err := service.UpdateTask(context.Background(), req)
	assert.Contains(t, err.Error(), "status invalid")
}

func TestGetTaskList(t *testing.T) {
	var (
		keys    = []string{"1", "2", "3"}
		expects = make([]*pbTask.Task, 3)
	)
	rmock.ExpectZRangeArgs(redis.ZRangeArgs{
		Key:    SortSet,
		ByLex:  true,
		Start:  "-",
		Stop:   "+",
		Offset: 0,
		Count:  30,
	}).SetVal(keys)

	for i, val := range keys {
		task := &pbTask.Task{
			Id:     val,
			Name:   fmt.Sprintf("test-%d", i),
			Status: 1,
		}
		data, _ := json.Marshal(task)
		expects[i] = task
		rmock.ExpectGet(fmt.Sprintf("%s:%s", TaskID, val)).SetVal(string(data))
	}

	resp, err := service.GetTaskList(ctx, &pbTask.GetTaskListRequest{PageSize: 30})
	assert.Nil(t, err)
	for i, expect := range expects {
		assert.Equal(t, expect, resp.Tasks[i])
	}
}

func TestGetTaskListWithToken(t *testing.T) {
	var (
		token   = utils.NewPageToken("25", 25)
		keys    = make([]string, 25)
		expects = make([]*pbTask.Task, 25)
	)
	for i := 0; i < len(keys); i++ {
		keys[i] = fmt.Sprintf("%d", i+26)
	}
	rmock.ExpectZRangeArgs(redis.ZRangeArgs{
		Key:    SortSet,
		ByLex:  true,
		Start:  "(" + token.GetID(),
		Stop:   "+",
		Offset: 0,
		Count:  25,
	}).SetVal(keys)

	for i, val := range keys {
		task := &pbTask.Task{
			Id:     val,
			Name:   fmt.Sprintf("test-%d", i),
			Status: 1,
		}
		data, _ := json.Marshal(task)
		expects[i] = task
		rmock.ExpectGet(fmt.Sprintf("%s:%s", TaskID, val)).SetVal(string(data))
	}

	resp, err := service.GetTaskList(ctx, &pbTask.GetTaskListRequest{PageToken: token.GetToken()})
	assert.Nil(t, err)
	for i, expect := range expects {
		assert.Equal(t, expect, resp.Tasks[i])
	}
	assert.NotEmpty(t, resp.NextToken)
}

// mock
type mockGenerator struct{}
type mockSequencer struct{}

func (m *mockGenerator) Next() (utils.Sequence, error) {
	return &mockSequencer{}, nil
}

func (m *mockSequencer) Uint64() uint64 {
	return num.Uint64()
}
func (m *mockSequencer) String() string {
	return num.String()
}
func (m *mockSequencer) Float64() float64 {
	result, _ := num.Float64()
	return result
}
