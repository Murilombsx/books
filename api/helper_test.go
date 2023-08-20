package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func executeHttp(t *testing.T, handler http.Handler, method string, url string) *httptest.ResponseRecorder {
	t.Helper()
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)
	handler.ServeHTTP(resp, req)
	return resp
}
