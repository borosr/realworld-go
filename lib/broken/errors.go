package broken

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	TypeValidation = "validation"
	TypeInternal   = "internal"
	TypeForbidden  = "forbidden"
)

type Mess interface {
	Error() string
	As(e any) bool
	Format() (string, int)
}

func New(typ, msg string) error {
	return &Thing{
		Type:    typ,
		Message: msg,
	}
}

type Thing struct {
	Message string `json:"message"`
	Type    string `json:"-"`
}

func (t *Thing) Format() (string, int) {
	marshal, _ := json.Marshal(t)
	var code int
	switch t.Type {
	case TypeValidation:
		code = http.StatusBadRequest
	case TypeForbidden:
		code = http.StatusForbidden
	case TypeInternal:
		code = http.StatusInternalServerError
	default:
		code = http.StatusInternalServerError
	}
	return string(marshal), code
}

func (t *Thing) Error() string {
	return t.Message
}

func (t *Thing) Is(err error) bool {
	return t.Message == err.Error()
}

func (t *Thing) As(e any) bool {
	switch et := e.(type) {
	case **Thing:
		*et = t
	case *Mess:
		*et = t
	case *error:
		*et = errors.New(t.Message)
	default:
		return false
	}
	return true
}

func Validation(msg string) error {
	return New(TypeValidation, msg)
}

func Forbidden(msg string) error {
	return New(TypeForbidden, msg)
}

func Forbiddenf(format string, args ...any) error {
	return New(TypeForbidden, fmt.Sprintf(format, args...))
}

func Internal(msg string) error {
	return New(TypeInternal, msg)
}
