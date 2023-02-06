// Conjunto de ferramentas para evitar boilerplates que reduzem a legibilidade do teste

package assert

import (
	"strings"
	"testing"
)

func StringContains(t *testing.T, s, substr string) {
	if !strings.Contains(s, substr) {
		t.Fatalf("assertion failed: '%s' does not contain '%s'", s, substr)
	}
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("assertion failed: expected no error but got '%s'", err.Error())
	}
}

func StatusCode(t *testing.T, got int, want int) {
	if want != got {
		t.Fatalf("assetion failed: expected status code '%d', but got '%d'", want, got)
	}
}
