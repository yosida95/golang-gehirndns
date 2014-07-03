package gehirndns

import (
	"errors"
)

var (
	ErrMaybeRegistered = errors.New("This record is maybe registered at Gehirn DNS.  Use `UpdateResource(IRecord) error` insted of this method")
	ErrIdUnset         = errors.New("Record id is unset")
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err errorResponse) Error() string {
	return err.Message
}
