package DAOHelper

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

type HttpError interface {
	error
	GetStatus() int
}

type DefaultHttpError struct {
	InnerError error
	StatusCode int
}

func (e *DefaultHttpError) Error() string {
	return e.InnerError.Error()
}

func (e *DefaultHttpError) GetStatus() int {
	if e.StatusCode < 100 {
		return http.StatusInternalServerError
	}

	return e.StatusCode
}

func ReportError(ctx context.Context, w http.ResponseWriter, e error) {
	if e == nil {
		return
	}

	httpError, isHttpError := e.(HttpError)
	if !isHttpError {
		httpError = &DefaultHttpError{InnerError: e}
	}

	log.Errorf(ctx, httpError.Error())
	http.Error(w, httpError.Error(), httpError.GetStatus())
}
