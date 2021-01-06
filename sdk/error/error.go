package error

import (
	"fmt"
)

type Err struct {
	Msg string
}

func NewError(msg string) Err {
	return Err{Msg: msg}
}

func (e Err) Error() string {
	return fmt.Sprintf("SDK error: %s", e.Msg)
}
