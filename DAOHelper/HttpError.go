package DAOHelper

import (
	"appengine"
	"net/http"
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

func ReportError(ctx appengine.Context, w http.ResponseWriter, e error) {
	if e == nil {
		return
	}

	httpError, isHttpError := e.(HttpError)
	if !isHttpError {
		httpError = &DefaultHttpError{InnerError: e}
	}

	ctx.Errorf(httpError.Error())
	http.Error(w, httpError.Error(), httpError.GetStatus())
}
