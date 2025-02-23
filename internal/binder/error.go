package binder

import "fmt"

type ErrNotFound struct {
	Type  string
	Value string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s variable {%s} not found in method params", e.Type, e.Value)
}

type ErrDuplicated struct {
	Type  string
	Value string
}

func (e *ErrDuplicated) Error() string {
	return fmt.Sprintf("duplicated %s variable {%s} found", e.Type, e.Value)
}
