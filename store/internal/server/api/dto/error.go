package dto

import (
	"net/http"
)

type ErrorData struct {
	RequestID  string `json:"request_id"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func InternalServerError(requestID, message string) *ErrorData {
	return &ErrorData{
		RequestID:  requestID,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func InvalidRequestError(requestID, message string) *ErrorData {
	return &ErrorData{
		RequestID:  requestID,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}
