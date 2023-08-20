package api

import (
	"net/http"
)

type readinessHandler struct {
}

func NewHeatlhHandler() http.Handler {
	return readinessHandler{}
}

func (h readinessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}
