package AppEngineHelper

import (
	"cloud.google.com/go/datastore"
)

// type MultiError []error
type MultiError struct {
	datastore.MultiError
}

func ToMultiError(err error) MultiError {
	if err == nil {
		return MultiError{}
	}

	if multiError, ok := err.(datastore.MultiError); ok {
		return MultiError{multiError}
	}

	return MultiError{[]error{err}}
}

func (errs MultiError) Filter(filterFn func(error, int) bool) MultiError {
	filteredErrors := make([]error, 0, len(errs.MultiError))
	for index, err := range errs.MultiError {
		if filterFn(err, index) {
			filteredErrors = append(filteredErrors, err)
		}
	}
	return MultiError{filteredErrors}
}

func (errs MultiError) ToError() error {
	if len(errs.MultiError) > 0 {
		return &errs
	}
	return nil
}
