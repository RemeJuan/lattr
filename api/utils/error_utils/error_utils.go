package error_utils

import (
	"encoding/json"
	"net/http"
)

type MessageErr interface {
	Message() string
	Status() int
	Error() string
}

type messageErr struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
}

func (e *messageErr) Error() string {
	return e.ErrError
}

func (e *messageErr) Message() string {
	return e.ErrMessage
}

func (e *messageErr) Status() int {
	return e.ErrStatus
}

func NotFoundError(message string) MessageErr {
	return &messageErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func ForbiddenError(message string) MessageErr {
	return &messageErr{
		ErrMessage: message,
		ErrStatus:  http.StatusForbidden,
		ErrError:   "bad_request",
	}
}
func UnprocessableEntityError(message string) MessageErr {
	return &messageErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "invalid_request",
	}
}

func ApiErrFromBytes(body []byte) (MessageErr, error) {
	var result messageErr
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func InternalServerError(message string) MessageErr {
	return &messageErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "server_error",
	}
}

func NotImplementedError(message string) MessageErr {
	return &messageErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotImplemented,
		ErrError:   "server_error",
	}
}
