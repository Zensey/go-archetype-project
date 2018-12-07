package utils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func emptyHandler(w http.ResponseWriter, r *http.Request) {}

func TestBasicAuth(t *testing.T) {
	var user = "user"
	var pass = "pass"
	hf := http.HandlerFunc(BasicAuth(emptyHandler, user, pass))

	req := httptest.NewRequest(http.MethodGet, "http://localhost:9999/test", nil)
	resp := httptest.NewRecorder()
	hf.ServeHTTP(resp, req)
	assert.Equal(t, 401, resp.Code)

	req.SetBasicAuth(user, pass)
	resp = httptest.NewRecorder()
	hf.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	hf = http.HandlerFunc(BasicAuth(emptyHandler, user, ""))
	req.SetBasicAuth(user, "pass~")
	resp = httptest.NewRecorder()
	hf.ServeHTTP(resp, req)
	assert.Equal(t, 401, resp.Code)

	hf = http.HandlerFunc(BasicAuth(emptyHandler, "", ""))
	req.SetBasicAuth("user~", pass)
	resp = httptest.NewRecorder()
	hf.ServeHTTP(resp, req)
	assert.Equal(t, 401, resp.Code)
}
