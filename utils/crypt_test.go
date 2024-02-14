package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCrypt(t *testing.T) {

	signingKey := "123456789012345678901234567890AB"
	originalText := "atirei o pau ao gato"

	cypher, err := Encrypt(originalText, signingKey)
	assert.NoError(t, err)

	decrypted, err := Decrypt(cypher, signingKey)
	assert.NoError(t, err)

	assert.Equal(t, originalText, decrypted)
}

func TestHashPassword(t *testing.T) {

	pass := "chuchas"
	hash, err := Hash(pass)
	assert.NoError(t, err)

	err = HashCompare(hash, pass)
	assert.NoError(t, err)

}
