package AppEngineHelper

import (
	"DAOHelper"
	"encoding/json"
	"net/http"
)

type HTTPHandler func(http.ResponseWriter, *http.Request) (interface{}, error)

func HandlerW(f HTTPHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dto, err := f(w, r)
		if err == nil && dto != nil {
			if dtoStr, ok := dto.(string); ok {
				w.Write([]byte(dtoStr))
			} else if dtoBytes, ok := dto.([]byte); ok {
				w.Write(dtoBytes)
			} else {
				_, err = WriteJSON(w, dto)
			}
		}

		if err != nil {
			DAOHelper.ReportError(ctx, w, err)
		}
	}
}

func WriteJSON(w http.ResponseWriter, data interface{}) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	if err == nil {
		w.Write(js)
	}

	return js, err
}

func CreateQueryParamString(params map[string]string) string {
	var paramStr string
	for k, v := range params {
		if len(paramStr) != 0 {
			paramStr = paramStr + "&"
		}
		paramStr = paramStr + k + "=" + v
	}
	return paramStr
}
