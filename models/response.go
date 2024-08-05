package models

import (
	"net/http"

	"sebi-scrapper/constants"

	"github.com/sinhashubham95/go-utils/errors"
)

// Response is the structure for the success or error response.
type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}

// GetErrorResponse is used to get the error response.
func GetErrorResponse(err error) (int, Response) {
	var r *errors.Error
	ok := errors.As(err, &r)
	if !ok {
		r = &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	return r.StatusCode, Response{Status: constants.ErrorAPIStatus, Data: r}
}

// GetNoContentSuccessResponse is used to get the no content response.
func GetNoContentSuccessResponse() (int, Response) {
	return http.StatusNoContent, Response{Status: constants.SuccessAPIStatus}
}

// GetOKSuccessResponse is used to get the response.
func GetOKSuccessResponse(data any) (int, Response) {
	return http.StatusOK, Response{Status: constants.SuccessAPIStatus, Data: data}
}

// GetNotFoundResponse is used to get not found response.
func GetNotFoundResponse() (int, Response) {
	return http.StatusNotFound, Response{Status: constants.ErrorAPIStatus}
}

// GetCreatedSuccessResponse is used to get status created response.
func GetCreatedSuccessResponse(data any) (int, Response) {
	return http.StatusCreated, Response{Status: constants.SuccessAPIStatus, Data: data}
}
