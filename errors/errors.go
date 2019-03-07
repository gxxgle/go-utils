package errors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gxxgle/go-utils/json"
	"github.com/gxxgle/go-utils/log"
	"github.com/micro/go-micro/errors"
)

var (
	Errors = make(map[int32]error)
)

func Parse(err error) *errors.Error {
	e := new(errors.Error)
	if errr := json.UnmarshalFromString(err.Error(), e); errr != nil {
		e.Detail = err.Error()
	}

	return e
}

func GetErrorCode(err error) int32 {
	e := Parse(err)
	return e.Code
}

func addError(err *errors.Error) *errors.Error {
	e, ok := Errors[err.Code]
	if ok {
		log.Fatalw("duplate error code", "exists_err", e, "new_err", err)
	}

	Errors[err.Code] = err
	return err
}

func BadRequest(code int32, detail string) error {
	return addError(&errors.Error{
		Code:   code,
		Detail: detail,
		Status: http.StatusText(400),
	})
}

func Unauthorized(code int32, detail string) error {
	return addError(&errors.Error{
		Code:   code,
		Detail: detail,
		Status: http.StatusText(401),
	})
}

func Forbidden(code int32, detail string) error {
	return addError(&errors.Error{
		Code:   code,
		Detail: detail,
		Status: http.StatusText(403),
	})
}

func NotFound(code int32, detail string) error {
	return addError(&errors.Error{
		Code:   code,
		Detail: detail,
		Status: http.StatusText(404),
	})
}

func Conflict(code int32, detail string) error {
	return addError(&errors.Error{
		Code:   code,
		Detail: detail,
		Status: http.StatusText(409),
	})
}

func Service(code int32, detail string) error {
	return addError(&errors.Error{
		Code:   code,
		Detail: detail,
		Status: http.StatusText(503),
	})
}

func Internal(format string, a ...interface{}) error {
	return &errors.Error{
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(500),
	}
}

// IsConnClosing return is grpc conn closing error
func IsConnClosing(err error) bool {
	return strings.Contains(err.Error(), "transport is closing")
}

// IsContextCanceled return is context canceled error
func IsContextCanceled(err error) bool {
	return strings.Contains(err.Error(), context.Canceled.Error())
}

// IsContextDeadlineExceeded return is context deadline exceeded error
func IsContextDeadlineExceeded(err error) bool {
	return strings.Contains(err.Error(), context.DeadlineExceeded.Error())
}

// CanStreamIgnoreError return is can stream ignore error
func CanStreamIgnoreError(err error) bool {
	return IsConnClosing(err) || IsContextCanceled(err) ||
		IsContextDeadlineExceeded(err)
}
