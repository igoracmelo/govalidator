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

// função para reduzir boilerplate nos testes de integração
// ela recebe um Reader que é a stream do body cria um request, um response writer,
// os passa para o meu handler e retorna só o objeto de response
func setupVerify(body io.Reader) *http.Response {
	req := httptest.NewRequest("POST", "/verify", body)
	w := httptest.NewRecorder()
	handler.HandleVerify(w, req)
	res := w.Result()
	return res
}

// sei que foi dito que a entrada vai ser sempre válida, mas acho importante tratar
// erros, principalmente coisas que derrubam o servidor como um nil pointer dereference
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

// testando uma única regra para garantir o comportamento correto
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

// conjunto de testes utilizando todas as regras de senha ao mesmo tempo para
// aumentar a cobertura de testes e validar cenários mais complexos
func TestAllRules(t *testing.T) {
	t.Run("should fail all rules", func(t *testing.T) {
		res := setupVerify(strings.NewReader(`{
			"password": "Ab1233",
			"rules": [
				{ "rule": "minSize",         "value": 7 },
				{ "rule": "minUppercase",    "value": 2 },
				{ "rule": "minLowercase",    "value": 3 },
				{ "rule": "minDigit",        "value": 5 },
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
				{ "rule": "minSize",         "value": 11 },
				{ "rule": "minUppercase",    "value": 2 },
				{ "rule": "minLowercase",    "value": 3 },
				{ "rule": "minDigit",        "value": 5 },
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

	t.Run("should fail some rules", func(t *testing.T) {
		res := setupVerify(strings.NewReader(`{
			"password": "AbabB12345@@",
			"rules": [
				{ "rule": "minSize",         "value": 11 },
				{ "rule": "minUppercase",    "value": 3 },
				{ "rule": "minLowercase",    "value": 3 },
				{ "rule": "minDigit",        "value": 6 },
				{ "rule": "minSpecialChars", "value": 2 },
				{ "rule": "noRepeted" }
			]
		}`))
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusOK)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, `"verify":false`)
		assert.StringContains(t, msg, `minUppercase`)
		assert.StringContains(t, msg, `minDigit`)
		assert.StringContains(t, msg, `noRepeted`)
	})
}

// testando caracteres especiais
func TestSpecialUnicodeCharacters(t *testing.T) {
	t.Run("should fail minSize because emoji ocupies a single size unit", func(t *testing.T) {
		res := setupVerify(strings.NewReader(`{
			"password": "😱",
			"rules": [
				{ "rule": "minSize", "value": 2 }
			]
		}`))
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusOK)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, `"verify":false`)
		assert.StringContains(t, msg, "minSize")
	})

	t.Run("should recognize uppercase and lowercase special chars", func(t *testing.T) {
		res := setupVerify(strings.NewReader(`{
			"password": "Áçãö",
			"rules": [
				{ "rule": "minSize", "value": 4 },
				{ "rule": "minUppercase", "value": 1 },
				{ "rule": "minLowercase", "value": 3 }
			]
		}`))
		defer res.Body.Close()
		assert.StatusCode(t, res.StatusCode, http.StatusOK)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		msg := string(data)
		assert.StringContains(t, msg, `"verify":true`)
	})
}
