// Package error предназначен для унифицированного пробороса ошибок между сервисами через gRPC.
package error

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
)

const (
	// ExternalGroup представляет внешние ошибки от ПС или других внешних сервисов.
	ExternalGroup = "external"
	// InternalGroup представляет внутренние ошибки наших сервисов.
	InternalGroup = "internal"
)

type Error struct {
	// Group внешняя или внутренняя ошибка
	Group string
	// Code унифицированный код ошибки в рамках всего сервиса
	Code string
	// Description описание ошибки
	Description string

	//internalCode внутренний код для работы с gRPC
	internalCode codes.Code
}

func NewExternal() Error {
	return Error{
		Group:        ExternalGroup,
		internalCode: codes.Unavailable,
	}
}

func NewInternal() Error {
	return Error{
		Group:        InternalGroup,
		internalCode: codes.Internal,
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
			Group:        detail.GetGroup(),
			Code:         detail.GetCode(),
			Description:  st.Message(),
			internalCode: code,
		}
	default:
		return nil
	}
}

func (e *Error) WithDescription(desc string) *Error {
	e.Description = desc

	return e
}

func (e *Error) WithCode(code string) *Error {
	e.Code = code

	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error occurred with code %q: %s", e.Group, e.Code, e.Description)
}

func (e *Error) GRPCStatus() *status.Status {
	details := pb.Error{
		Group: e.Group,
		Code:  e.Code,
	}

	st, err := status.New(e.internalCode, e.Description).WithDetails(&details)
	if err != nil {
		return status.New(codes.Unknown, fmt.Sprintf("failed to pack details for error code %q", e.Code))
	}

	return st
}
