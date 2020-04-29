package errors

import (
	"fmt"

	"github.com/gxxgle/go-utils/json"
	"github.com/gxxgle/go-utils/log"
)

var (
	UnrecognizedErrorCode = int32(1)
	InternalErrorCode     = int32(2)

	allErrors = map[int32]error{
		UnrecognizedErrorCode: &Error{Code: UnrecognizedErrorCode, Detail: "unrecognized error"},
		InternalErrorCode:     &Error{Code: InternalErrorCode, Detail: "internal error"},
	}
)

type Error struct {
	Code     int32  `json:"code"`
	Detail   string `json:"detail"`
	Internal string `json:"internal,omitempty"`
}

func (e *Error) Error() string {
	return json.MustMarshalToString(e)
}

func Parse(err error) *Error {
	e := new(Error)
	if errr := json.UnmarshalFromString(err.Error(), e); errr != nil {
		e.Code = UnrecognizedErrorCode
		e.Detail = err.Error()
	}

	return e
}

func GetErrorCode(err error) int32 {
	return Parse(err).Code
}

func addError(err *Error) *Error {
	e, ok := allErrors[err.Code]
	if ok {
		log.Fatalw("duplate error code", "exists_err", e, "new_err", err)
	}

	allErrors[err.Code] = err
	return err
}

func New(code int32, detail string) error {
	return addError(&Error{
		Code:   code,
		Detail: detail,
	})
}

func Internal(format string, a ...interface{}) error {
	return &Error{
		Code:     InternalErrorCode,
		Detail:   "internal error",
		Internal: fmt.Sprintf(format, a...),
	}
}
