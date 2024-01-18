package helper

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NotFoundErr(msg string, field string, resourceId string) error {
	st := status.New(codes.NotFound, msg)
	v := &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: fmt.Sprintf("'%s' was not found", resourceId),
	}

	badReq := &errdetails.BadRequest{}
	badReq.FieldViolations = append(badReq.FieldViolations, v)

	st, _ = st.WithDetails(badReq)
	return st.Err()
}

func RequiredFieldErr(msg string, field string) error {
	st := status.New(codes.InvalidArgument, msg)
	v := &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: fmt.Sprintf("missing field '%s' is required", field),
	}

	badReq := &errdetails.BadRequest{}
	badReq.FieldViolations = append(badReq.FieldViolations, v)

	st, _ = st.WithDetails(badReq)
	return st.Err()
}

func InvalidErr(msg string, field string, value interface{}) error {
	st := status.New(codes.InvalidArgument, msg)
	v := &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: fmt.Sprintf("value '%v' is invalid", value),
	}

	badReq := &errdetails.BadRequest{}
	badReq.FieldViolations = append(badReq.FieldViolations, v)

	st, _ = st.WithDetails(badReq)
	return st.Err()
}

func BadRequestErr(msg string, field string, description string) error {
	st := status.New(codes.InvalidArgument, msg)
	v := &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: description,
	}

	badReq := &errdetails.BadRequest{}
	badReq.FieldViolations = append(badReq.FieldViolations, v)

	st, _ = st.WithDetails(badReq)
	return st.Err()
}

func InternalErr(msg string) error {
	st := status.New(codes.Internal, msg)
	return st.Err()
}
