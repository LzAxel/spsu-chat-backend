package apperror

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

type AppError struct {
	Err     error
	Service string
	Func    string
	Data    map[string]interface{}
}

func (e AppError) Error() string {
	return fmt.Sprintf(
		"app error:%s.%s:%s:data=%v",
		e.Service, e.Func,
		e.Err.Error(),
		e.Data,
	)
}

func NewAppError(err error, service, funcName string, data map[string]interface{}) AppError {
	return AppError{
		Err:     err,
		Service: service,
		Func:    funcName,
		Data:    data,
	}
}

type DBError struct {
	Err     error
	Service string
	Func    string

	Query string
	Args  []interface{}
}

func (e DBError) Error() string {
	return fmt.Sprintf(
		"db error:%s.%s:%s:query=%s args=%v",
		e.Service, e.Func,
		e.Err.Error(),
		e.Query,
		e.Args,
	)
}

func NewDBError(err error, service, funcName string, query string, args []interface{}) DBError {
	return DBError{
		Err:     err,
		Service: service,
		Func:    funcName,
		Query:   query,
		Args:    args,
	}
}
