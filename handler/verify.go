package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

type VerifyRequest struct {
	Password string `json:"password"`
	Rules    []struct {
		Rule  string `json:"rule"`
		Value int    `json:"value"`
	}
}

type VerifyResponse struct {
	Verify  bool     `json:"verify"`
	NoMatch []string `json:"noMatch"`
}

func HandleVerify(w http.ResponseWriter, r *http.Request) {
	var body VerifyRequest
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		msg := fmt.Sprintf("Failed to parse request body: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(body.Password) == "" {
		http.Error(w, "required field 'password' missing", http.StatusBadRequest)
		return
	}

	// if len(body.Rules) == 0 {
	// 	http.Error(w, "no rule specified", http.StatusBadRequest)
	// 	return
	// }

	res := VerifyResponse{
		Verify:  true,
		NoMatch: []string{},
	}

	for _, rule := range body.Rules {
		validate, ok := PasswordValidators[rule.Rule]

		if !ok {
			http.Error(w, fmt.Sprintf("Unkown rule '%s'", rule.Rule), http.StatusBadRequest)
			return
		}

		if !validate(body.Password, rule.Value) {
			res.Verify = false
			res.NoMatch = append(res.NoMatch, rule.Rule)
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate response body: %s", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

type PasswordValidator func(string, int) bool

var PasswordValidators = map[string]PasswordValidator{
	"minSize": func(pass string, x int) bool {
		return len(pass) >= x
	},

	"minUppercase": func(pass string, x int) bool {
		count := 0
		for _, r := range []rune(pass) {
			if unicode.IsUpper(r) {
				count++
			}
		}
		return count >= x
	},

	"minLowercase": func(pass string, x int) bool {
		count := 0
		for _, r := range []rune(pass) {
			if unicode.IsLower(r) {
				count++
			}
		}
		return count >= x
	},

	"minDigit": func(pass string, x int) bool {
		count := 0
		for _, r := range []rune(pass) {
			if unicode.IsDigit(r) {
				count++
			}
		}
		return count >= x
	},

	"minSpecialChars": func(pass string, x int) bool {
		specialChars := "!@#$%^&*()-+\\/{}[]"
		count := 0
		for _, r := range []rune(pass) {
			if strings.ContainsRune(specialChars, r) {
				count++
			}
		}
		return count >= x
	},

	// este nome está com erro de digitação no PDF, mas mantive caso sejam utilizadas
	// ferramentas de testes automatizados que se baseiem no nome que está no PDF
	"noRepeted": func(pass string, _ int) bool {
		runes := []rune(pass)
		for i := range runes[1:] {
			if runes[i] == runes[i-1] {
				return false
			}
		}
		return true
	},
}
