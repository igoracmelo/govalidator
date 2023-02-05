package handler_test

import (
	"estudiosol/assert"
	"estudiosol/handler"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupVerify(body io.Reader) (res *http.Response) {
	req := httptest.NewRequest("POST", "/verify", body)
	w := httptest.NewRecorder()
	handler.HandleVerify(w, req)
	res = w.Result()
	return
}

func TestInvalidRequestBody(t *testing.T) {
	t.Run("nil body should return bad request", func(t *testing.T) {
		res := setupVerify(nil)
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusBadRequest)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, "parse request")
	})

	t.Run("empty password should return bad request", func(t *testing.T) {
		res := setupVerify(strings.NewReader("{}"))
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
		res := setupVerify(strings.NewReader(`{"password": "123"}`))
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusOK)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, `{"verify":true,"noMatch":[]}`)
	})
}

func TestAllRules(t *testing.T) {
	t.Run("should fail all rules", func(t *testing.T) {
		res := setupVerify(strings.NewReader(`{
			"password": "Ab1233",
			"rules": [
				{ "rule": "minSize", "value": 7 },
				{ "rule": "minUppercase", "value": 2 },
				{ "rule": "minLowercase", "value": 3 },
				{ "rule": "minDigit", "value": 5 },
				{ "rule": "minSpecialChars", "value": 1 },
				{ "rule": "noRepeted" }
			]
		}`))
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusOK)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, `"verify":false`)
		assert.StringContains(t, msg, `minSize`)
		assert.StringContains(t, msg, `minUppercase`)
		assert.StringContains(t, msg, `minLowercase`)
		assert.StringContains(t, msg, `minDigit`)
		assert.StringContains(t, msg, `minSpecialChars`)
		assert.StringContains(t, msg, `noRepeted`)
	})

	t.Run("should pass all rules", func(t *testing.T) {
		res := setupVerify(strings.NewReader(`{
			"password": "AbabB12345@",
			"rules": [
				{ "rule": "minSize", "value": 11 },
				{ "rule": "minUppercase", "value": 2 },
				{ "rule": "minLowercase", "value": 3 },
				{ "rule": "minDigit", "value": 5 },
				{ "rule": "minSpecialChars", "value": 1 },
				{ "rule": "noRepeted" }
			]
		}`))
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusOK)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, `{"verify":true,"noMatch":[]}`)
	})
}

func TestMinSize(t *testing.T) {
	t.Run("size less than minSize should not be accepted", func(t *testing.T) {
		res := setupVerify(strings.NewReader(`{
			"password": "123",
			"rules": [{
				"rule": "minSize",
				"value": 4
			}]
		}`))
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusOK)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, `{"verify":false,"noMatch":["minSize"]}`)
	})
}

func TestInvalidRule(t *testing.T) {
	t.Run("invalid rule should return bad request", func(t *testing.T) {
		res := setupVerify(strings.NewReader(`{
			"password": "abc123#",
			"rules": [{
				"rule": "invalid rule",
				"value": 4
			}]
		}`))
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusBadRequest)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, "Unknown rule")
	})
}
