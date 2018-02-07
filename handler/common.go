package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Meta    interface{} `json:"meta"`
}

func newJson(w http.ResponseWriter, header int, res response) []byte {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
	r, err := json.Marshal(&res)
	if err != nil {
		log.Printf("%v", res)
		log.Fatalf("could not marshal json responce: %s", err)
	}
	return r
}

func success(w http.ResponseWriter, data interface{}, meta interface{}) {
	res := newJson(w, http.StatusOK, response{Success: true, Data: data, Meta: meta})
	w.Write(res)
}

func successCreated(w http.ResponseWriter, data interface{}) {
	res := newJson(w, http.StatusCreated, response{Success: true, Data: data})
	w.Write(res)
}

func serverError(w http.ResponseWriter, err error) {
	log.Printf("%v", err)
	res := newJson(w, http.StatusInternalServerError, response{Error: err.Error()})
	w.Write(res)
}

func badRequest(w http.ResponseWriter, err string) {
	res := newJson(w, http.StatusBadRequest, response{Error: err})
	w.Write(res)
}

func queryParamInt(r *http.Request, name string, defaultValue int) (int, error) {
	v := r.FormValue(name)
	if v == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(v)
}
