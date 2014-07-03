package gehirndns

import (
	"errors"
)

var (
	ErrMaybeRegistered = errors.New("This record is maybe registered at Gehirn DNS.  Use `UpdateResource(IRecord) error` insted of this method")
)
