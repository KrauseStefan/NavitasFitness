package DAOHelper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorType string

const (
	UniqueConstraint ErrorType = "unique_constraint"
	Invalid          ErrorType = "invalid"
)

type ConstraintError struct {
	Field   string    `json:"field"`
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
}

func (e ConstraintError) GetStatus() int {
	switch e.Type {
	case UniqueConstraint:
		return http.StatusConflict
	case Invalid:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (e ConstraintError) Error() string {
	if len(e.Message) == 0 {
		if e.Type == UniqueConstraint {
			e.Message = fmt.Sprintf("Cannot create user, %s already in use", e.Field)
		} else if e.Type == Invalid {
			e.Message = fmt.Sprintf("Cannot create user, %s is invalid", e.Field)
		}
	}

	js, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return string(js)
}
