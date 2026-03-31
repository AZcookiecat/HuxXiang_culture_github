package app

import (
	"errors"
	"fmt"
	"net/http"
)

type StatusError interface {
	error
	StatusCode() int
	PublicMessage() string
}

type AppError struct {
	Status  int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) StatusCode() int {
	if e.Status <= 0 {
		return http.StatusInternalServerError
	}
	return e.Status
}

func (e *AppError) PublicMessage() string {
	if e.Message == "" {
		return "服务器开小差了，请稍后重试"
	}
	return e.Message
}

func NewError(status int, message string, err error) *AppError {
	return &AppError{Status: status, Message: message, Err: err}
}

func ResolveError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var statusErr StatusError
	if errors.As(err, &statusErr) {
		return statusErr.StatusCode(), statusErr.PublicMessage()
	}

	return http.StatusInternalServerError, "服务器开小差了，请稍后重试"
}
