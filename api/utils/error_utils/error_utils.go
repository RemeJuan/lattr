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

type MessageErrStruct struct {
	ErrMessage string `json:"message" example:"Invalid body"`
	ErrStatus  int    `json:"status" example:"400"`
	ErrError   string `json:"error" example:"bad_request"`
}

func (e *MessageErrStruct) Error() string {
	return e.ErrError
}

func (e *MessageErrStruct) Message() string {
	return e.ErrMessage
}

func (e *MessageErrStruct) Status() int {
	return e.ErrStatus
}

func NotFoundError(message string) MessageErr {
	return &MessageErrStruct{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func ForbiddenError(message string) MessageErr {
	return &MessageErrStruct{
		ErrMessage: message,
		ErrStatus:  http.StatusForbidden,
		ErrError:   "bad_request",
	}
}
func UnprocessableEntityError(message string) MessageErr {
	return &MessageErrStruct{
		ErrMessage: message,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "invalid_request",
	}
}

func InternalServerError(message string) MessageErr {
	return &MessageErrStruct{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "server_error",
	}
}

func NotImplementedError(message string) MessageErr {
	return &MessageErrStruct{
		ErrMessage: message,
		ErrStatus:  http.StatusNotImplemented,
		ErrError:   "server_error",
	}
}

func ApiErrFromBytes(body []byte) (MessageErr, error) {
	var result MessageErrStruct
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
