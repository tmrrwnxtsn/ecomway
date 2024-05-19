// Package error предназначен для унифицированного пробороса ошибок между сервисами через gRPC.
package error

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
)

type Group string

const (
	// GroupExternal представляет внешние ошибки от ПС или других внешних сервисов.
	GroupExternal Group = "external"
	// GroupInternal представляет внутренние ошибки наших сервисов.
	GroupInternal Group = "internal"
)

type Code string

const (
	// CodeObjectNotFound используется, когда искомый объект, для изменения или удаления, не найден.
	CodeObjectNotFound Code = "object not found"
	// CodeUnresolvedStatusConflict используется, когда действие не может быть соверешно над объектом из-за неподходящего состояния.
	CodeUnresolvedStatusConflict Code = "unresolved status for action"
)

type Error struct {
	Group       Group
	Code        Code
	Description string

	grpcCode codes.Code
}

func NewExternal() *Error {
	return &Error{
		Group:    GroupExternal,
		grpcCode: codes.Unavailable,
	}
}

func NewInternal() *Error {
	return &Error{
		Group:    GroupInternal,
		grpcCode: codes.Internal,
	}
}

func FromProto(err error) *Error {
	st := status.Convert(err)

	code := st.Code()
	if code != codes.Unavailable && code != codes.Internal {
		return nil
	}

	details := st.Details()
	if len(details) == 0 {
		return nil
	}

	switch detail := details[0].(type) { // ожидаем только один элемент в массиве
	case *pb.Error:
		return &Error{
			Group:       Group(detail.GetGroup()),
			Code:        Code(detail.GetCode()),
			Description: st.Message(),
			grpcCode:    code,
		}
	default:
		return nil
	}
}

func (e *Error) WithDescription(desc string) *Error {
	e.Description = desc

	return e
}

func (e *Error) WithCode(code Code) *Error {
	e.Code = code

	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error occurred with code %q: %s", e.Group, e.Code, e.Description)
}

func (e *Error) GRPCStatus() *status.Status {
	details := pb.Error{
		Group: string(e.Group),
		Code:  string(e.Code),
	}

	st, err := status.New(e.grpcCode, e.Description).WithDetails(&details)
	if err != nil {
		return status.New(codes.Unknown, fmt.Sprintf("failed to pack details for error code %q", e.Code))
	}

	return st
}
