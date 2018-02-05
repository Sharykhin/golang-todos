package handler

import (
	"log"
	"net/http"

	"github.com/sharykhin/todoapp/service/response"
)

type Handler struct {
}

func (h Handler) success(w http.ResponseWriter, data interface{}) {
	res, _ := response.NewJson(w, http.StatusOK, response.Response{true, data, nil})
	w.Write(res)
}

func (h Handler) successCreated(w http.ResponseWriter, data interface{}) {
	res, _ := response.NewJson(w, http.StatusCreated, response.Response{true, data, nil})
	w.Write(res)
}

func (h Handler) serverError(w http.ResponseWriter, err error) {
	log.Println(err)
	res, _ := response.NewJson(w, http.StatusInternalServerError, response.Response{false, nil, err.Error()})
	w.Write(res)
}

func (h Handler) queryParam(r *http.Request, name string, defaultValue string) string {
	v := r.FormValue(name)
	if v == "" {
		return defaultValue
	}
	return v
}
