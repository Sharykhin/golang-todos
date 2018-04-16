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

type successResponse struct {
	response
	Error error `json:"error"`
}

func newJSON(w http.ResponseWriter, header int, res response) []byte {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
	var r []byte
	var err error
	if res.Success == true {
		successR := successResponse{response: res, Error: nil}
		r, err = json.Marshal(&successR)
	} else {
		r, err = json.Marshal(&res)
	}

	if err != nil {
		log.Printf("%v", res)
		log.Fatalf("could not marshal json response: %s", err)
	}
	return r
}

func success(w http.ResponseWriter, data interface{}, meta interface{}) {
	res := newJSON(w, http.StatusOK, response{Success: true, Data: data, Meta: meta})
	_, err := w.Write(res)
	if err != nil {
		log.Fatalf("Could not sent a response: %v", err)
	}
}

func successCreated(w http.ResponseWriter, data interface{}) {
	res := newJSON(w, http.StatusCreated, response{Success: true, Data: data})
	w.Write(res)
}

func serverError(w http.ResponseWriter, err error) {
	log.Printf("%v", err)
	res := newJSON(w, http.StatusInternalServerError, response{Error: err.Error()})
	w.Write(res)
}

func badRequest(w http.ResponseWriter, err string) {
	res := newJSON(w, http.StatusBadRequest, response{Error: err})
	w.Write(res)
}

func queryParamInt(r *http.Request, name string, defaultValue int) (int, error) {
	v := r.FormValue(name)
	if v == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(v)
}
