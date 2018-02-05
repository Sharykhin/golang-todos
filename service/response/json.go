package response

import (
	"encoding/json"
	"net/http"
)

func NewJson(w http.ResponseWriter, header int, res Response) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
	return json.Marshal(&res)
}
