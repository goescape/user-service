package fault

import (
	"fmt"
	"log"
	"net/http"
	"user-svc/helpers/response"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrBadRequest     ErrorCode = "BAD_REQUEST"
	ErrTimeout        ErrorCode = "TIMEOUT"
	ErrConflict       ErrorCode = "CONFLICT"
	ErrUnprocessable  ErrorCode = "UNPROCESSABLE_ENTITY"
	ErrForbidden      ErrorCode = "FORBIDDEN"
	ErrUnknown        ErrorCode = "UNKNOWN"
	ErrUnavailable    ErrorCode = "UNAVAILABLE"
)

type errorMessage string

const (
	msgInternalServer errorMessage = "An error occurred on the server. Please try again later."
	msgUnauthorized   errorMessage = "You are not authorized to perform this action."
	msgNotFound       errorMessage = "The requested data was not found."
	msgBadRequest     errorMessage = "Invalid request. Please check the submitted data."
	msgTimeout        errorMessage = "The request timed out. Please try again."
	msgConflict       errorMessage = "The submitted data already exists or there is a conflict."
	msgUnprocessable  errorMessage = "The request could not be processed."
	msgForbidden      errorMessage = "You're not in the right place!"
	msgUnknown        errorMessage = "An unknown error occurred."
	msgUnavailable    errorMessage = "Service unavailable"
)

var errorMessages = map[ErrorCode]errorMessage{
	ErrInternalServer: msgInternalServer,
	ErrUnauthorized:   msgUnauthorized,
	ErrNotFound:       msgNotFound,
	ErrBadRequest:     msgBadRequest,
	ErrTimeout:        msgTimeout,
	ErrConflict:       msgConflict,
	ErrUnprocessable:  msgUnprocessable,
	ErrForbidden:      msgForbidden,
	ErrUnavailable:    msgUnavailable,
}

type ErrorResponse struct {
	HTTPStatus int    `json:"http_status"`
	Message    string `json:"message"`
}

type DetailedError struct {
	External ErrorResponse `json:"external"`
	Internal ErrorResponse `json:"internal"`
}

func GetExternalMessage(code ErrorCode) string {
	if msg, ok := errorMessages[code]; ok {
		return string(msg)
	}
	return string(msgUnknown)
}

func (e *DetailedError) Error() string {
	return fmt.Sprintf("External: %s | Internal: %s", e.External.Message, e.Internal.Message)
}

func newError(httpStatus int, code ErrorCode, internalMessage string) *DetailedError {
	return &DetailedError{
		External: ErrorResponse{
			HTTPStatus: httpStatus,
			Message:    GetExternalMessage(code),
		},
		Internal: ErrorResponse{
			HTTPStatus: httpStatus,
			Message:    internalMessage,
		},
	}
}

func Custom(httpStatus int, code ErrorCode, internalMessage string) *DetailedError {
	return newError(httpStatus, code, internalMessage)
}

func Response(ctx *gin.Context, err error) {
	errors, ok := err.(*DetailedError)
	if !ok {
		errors = newError(http.StatusInternalServerError, "Something went wrong", err.Error())
	}

	if errors.External.HTTPStatus >= http.StatusUnauthorized {
		log.Printf("[ERROR] status=%d | %s\n",
			errors.Internal.HTTPStatus,
			errors.Internal.Message,
		)
	}

	response.JSON(ctx, errors.External.HTTPStatus, errors.External.Message, nil)
}
