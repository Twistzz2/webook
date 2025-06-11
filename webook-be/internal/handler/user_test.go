package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestEncrypt(t *testing.T) {
	password := "hello@world123"
	t.Logf("Original password: %s", password)

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to encrypt password: %v", err)
	}
	t.Logf("Encrypted password: %s", string(encryptedPassword))

	// 1. 测试加密
	err = bcrypt.CompareHashAndPassword(encryptedPassword, []byte(password))
	assert.NoError(t, err)
}
