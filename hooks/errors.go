package hooks

import (
	"net/http"
)

type HookError struct {
	err string
}

func (e *HookError) Error() string {
	return e.err
}

func (e *HookError) Code() int {
	return http.StatusInternalServerError // 500
}
