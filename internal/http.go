package internal

import (
	"encoding/json"
	"net/http"
)

// WriteJSON marhsals data into json format using encoding/json and writes it to
// the reposnse writer
//
// This sets status code to code and Content-Type to JSON
func WriteJSON(w http.ResponseWriter, data interface{}, code int) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	return err
}
