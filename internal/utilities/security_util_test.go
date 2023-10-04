package utilities

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateJwtToken(t *testing.T) {
	userId := 101
	roles := []string{"USER"}
	_, err := GenerateJwtToken(userId, roles, (time.Second * 80))
	if err != nil {
		t.Error("Failed to generate token")
	}
}

func TestGenerateOpaqueToken(t *testing.T) {
	randomCharsLength := 10
	token := GenerateOpaqueToken(randomCharsLength)
	if token == "" {
		t.Error("Failed to generate Opaque token")
	}
}

func BenchmarkGenerateJwtToken(b *testing.B) {
	userId := 101
	roles := []string{"ADMIN"}
	for i := 0; i < b.N; i++ {
		_, err := GenerateJwtToken(userId, roles, (time.Second * 600))
		if err != nil {
			b.Error("Failed to generate Token")
		}
	}
}

func TestValidateJwtAndGetClaims(t *testing.T) {
	userId := 102
	roles := []string{"USER"}
	token, err := GenerateJwtToken(userId, roles, (time.Second * 80))
	if err != nil {
		t.Error("Failed to generate token")
	}
	_, err = ValidateJwtAndGetClaims(token)
	if err != nil {
		t.Errorf("%s token is invalid", token)
	}

}

func TestStrongPasswordCheck(t *testing.T) {
	var tests = []struct {
		password string
		want     bool
	}{
		{"1234", false},
		{"passwords", false},
		{"W_P12xxx10202@", true},
		{"073148291", false},
		{"@@@@@@AAA3333444511122abbb_000", true},
		{"..,a,a,,s,s.s11122222AAAa", true},
		{"...a.a.aa.a.a.a..a.a   wwwAAA221", true},
		{"01001029299292092AAAAAAAA2", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.password)
		t.Run(testname, func(t *testing.T) {
			answer := StrongPasswordCheck(tt.password)
			if answer != tt.want {
				t.Error("Expected ", tt.want)
			}

		})
	}

}
