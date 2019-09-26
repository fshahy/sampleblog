/*
Package helper implements the response sent back to clients.
It has the following fields:
	Status
	Message
	Data
*/
package helper

import (
	"errors"
	"reflect"

	"github.com/fshahy/sampleblog/article"
)

// Response represents responses sent back to clients.
type Response struct {
	Status  *int        `json:"status,omitempty"`
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

// NewResponse creates a new response using parameters passed to it.
func NewResponse(status int, message string, data interface{}) (Response, error) {
	if reflect.TypeOf(data).Kind() == reflect.Int64 {
		dt, ok := data.(int64)
		if !ok {
			err := errors.New("invalid data in creating response, expecting int64")

			return Response{Status: &status, Message: &message}, err
		}

		artcl := article.Article{
			ID: &dt,
		}

		return Response{
			Status:  &status,
			Message: &message,
			Data:    artcl,
		}, nil
	}

	dt, ok := data.([]article.Article)
	if !ok {
		err := errors.New("invalid data in creating response, expecting a slice of structs")

		return Response{Status: &status, Message: &message}, err
	}

	return Response{
		Status:  &status,
		Message: &message,
		Data:    &dt,
	}, nil
}
