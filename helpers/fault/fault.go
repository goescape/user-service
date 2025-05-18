package fault

import (
	"fmt"
	"log"
	"net/http"
	"user-svc/helpers/response"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

// Kumpulan constant kode error yang digunakan secara konsisten dalam aplikasi.

const (
	// ErrInternalServer Digunakan saat terjadi error di sisi server (500)
	ErrInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	// ErrUnauthorized Untuk error otorisasi pengguna (401)
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	// ErrNotFound Data tidak ditemukan (404)
	ErrNotFound ErrorCode = "NOT_FOUND"
	// ErrBadRequest Request tidak valid (400)
	ErrBadRequest ErrorCode = "BAD_REQUEST"
	// ErrTimeout Waktu request habis (timeout)
	ErrTimeout ErrorCode = "TIMEOUT"
	// ErrConflict Konflik data, misal duplikat (409)
	ErrConflict ErrorCode = "CONFLICT"
	// ErrUnprocessable Server paham request tapi tidak bisa diproses (422)
	ErrUnprocessable ErrorCode = "UNPROCESSABLE_ENTITY"
	// ErrForbidden Akses ditolak meskipun terautentikasi (403)
	ErrForbidden ErrorCode = "FORBIDDEN"
	// ErrUnknown Error tidak diketahui
	ErrUnknown ErrorCode = "UNKNOWN"
	// ErrUnavailable Service sedang tidak tersedia (503)
	ErrUnavailable ErrorCode = "UNAVAILABLE"
)

type errorMessage string

// Pesan user-friendly untuk masing-masing ErrorCode.

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

// Pemetaan kode error ke pesan eksternal
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
	HTTPStatus int    `json:"http_status"` // Status HTTP untuk response
	Message    string `json:"message"`     // Pesan kesalahan
}

type DetailedError struct {
	External ErrorResponse `json:"external"` // Pesan untuk user
	Internal ErrorResponse `json:"internal"` // Pesan untuk log/internal
}

func GetExternalMessage(code ErrorCode) string {
	// Ambil pesan user-friendly berdasarkan ErrorCode
	if msg, ok := errorMessages[code]; ok {
		return string(msg)
	}
	return string(msgUnknown)
}

func (e *DetailedError) Error() string {
	// Format error yang bisa digunakan sebagai string
	return fmt.Sprintf("External: %s | Internal: %s", e.External.Message, e.Internal.Message)
}

func newError(httpStatus int, code ErrorCode, internalMessage string) *DetailedError {
	// Helper untuk buat error detail
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
	// Fungsi utama untuk buat custom error
	return newError(httpStatus, code, internalMessage)
}

func Response(ctx *gin.Context, err error) {
	// Handler untuk mengirim response error ke client
	errors, ok := err.(*DetailedError)
	if !ok {
		// Fallback jika error bukan tipe DetailedError
		errors = newError(http.StatusInternalServerError, "Something went wrong", err.Error())
	}

	if errors.External.HTTPStatus >= http.StatusUnauthorized {
		// Log hanya untuk error mulai dari status 401 ke atas
		log.Printf("[ERROR] status=%d | %s\n",
			errors.Internal.HTTPStatus,
			errors.Internal.Message,
		)
	}

	// Kirim response JSON ke client
	response.JSON(ctx, errors.External.HTTPStatus, errors.External.Message, nil)
}
