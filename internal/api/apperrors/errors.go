package apperrors

import "net/http"

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(code, message string) *AppError {
	return &AppError{Code: code, Message: message, HTTPStatus: http.StatusBadRequest}
}

func NotFound(message string) *AppError {
	return &AppError{Code: "not_found", Message: message, HTTPStatus: http.StatusNotFound}
}

func Internal(message string) *AppError {
	return &AppError{Code: "internal", Message: message, HTTPStatus: http.StatusInternalServerError}
}
