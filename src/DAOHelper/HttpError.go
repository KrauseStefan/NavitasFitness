package DAOHelper

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
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
	if errs, ok := e.InnerError.(appengine.MultiError); ok && len(errs) > 1 {
		strs := make([]string, len(errs)+1)
		strs[0] = "Multi Error:"
		for i, err := range errs {
			if err != nil {
				strs[i+1] = err.Error()
			}
		}
		return strings.Join(strs, "\n")
	} else {
		return e.InnerError.Error()
	}
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
