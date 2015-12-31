package Services
import (
	"net/http"
	"encoding/json"
)

func writeJSON(w http.ResponseWriter, data interface{}) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(data)

	w.Write(js)

	return js, err
}