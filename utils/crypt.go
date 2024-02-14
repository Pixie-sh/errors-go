package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// Encrypt string. key must be 32bits for AES256
func Encrypt(stringToEncrypt string, keyString string) (string, error) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt string. key must be 32bits for AES256
func Decrypt(encryptedString string, keyString string) (string, error) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}

// Hash hash plain text. only the first cost is used
func Hash(originalText string, cost ...int) (string, error) {
	hashCost := bcrypt.DefaultCost
	if len(cost) > 0 {
		hashCost = cost[0]
	}
	blob, err := bcrypt.GenerateFromPassword([]byte(originalText), hashCost)
	return string(blob), err
}

// HashCompare compare plaint text with a hash
func HashCompare(hash string, plainText string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainText))
}

// Md5 create md5 sum
func Md5(keys ...string) string {
	h := md5.New()
	for _, key := range keys {
		_, _ = io.WriteString(h, key)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
