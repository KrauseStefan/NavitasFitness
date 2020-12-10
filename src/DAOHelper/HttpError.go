package DAOHelper

import (
	"net/http"
	"strings"

	log "github.com/KrauseStefan/NavitasFitness/logger"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

type HttpError interface {
	error
	GetStatus() int
}

type DefaultHttpError struct {
	InnerError error
	StatusCode int
}

func extractMultiErrors(multiError datastore.MultiError) string {
	strs := make([]string, len(multiError)+1)
	strs[0] = "Multi Error:"
	for i, err := range multiError {
		if err == nil {
			strs[i+1] = "nil"
		} else {
			strs[i+1] = err.Error()
		}
	}
	return strings.Join(strs, "\n")
}

func (e *DefaultHttpError) Error() string {
	if multiError, ok := e.InnerError.(datastore.MultiError); ok && len(multiError) > 1 {
		return extractMultiErrors(multiError)
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
