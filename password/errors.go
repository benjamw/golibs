package password

import (
	"fmt"
	"net/http"
)

type InvalidPasswordError struct {
	Msg string
}

func (e *InvalidPasswordError) Error() string {
	return fmt.Sprintf("The password submitted is invalid: %s", e.Msg)
}

func (e *InvalidPasswordError) Code() int {
	return http.StatusBadRequest
}
