package main

import (
	"estudiosol/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidRequestBody(t *testing.T) {
	t.Run("nil body should return bad request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/verify", nil)
		w := httptest.NewRecorder()

		handlerVerify(w, req)
		res := w.Result()
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusBadRequest)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, "parse request")
	})

	t.Run("empty password should return bad request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/verify", strings.NewReader("{}"))
		w := httptest.NewRecorder()

		handlerVerify(w, req)
		res := w.Result()
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusBadRequest)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, "password")
		assert.StringContains(t, msg, "missing")
	})

}

// fiquei na dúvida se deveria dar bad request ou permitir rules vazio,
// mas como o client que decide quais regras quer passar, optei por permitir
func TestEmptyRules(t *testing.T) {
	t.Run("empty 'rules' should be permitted", func(t *testing.T) {
		body := strings.NewReader(`{"password": "123"}`)
		req := httptest.NewRequest("POST", "/verify", body)
		w := httptest.NewRecorder()

		handlerVerify(w, req)
		res := w.Result()
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusOK)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, `{"verify":true,"noMatch":[]}`)
	})
}
