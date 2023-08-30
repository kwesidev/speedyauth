package utilities

import (
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
