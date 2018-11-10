package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt is encrypt
func Encrypt(data string, privateKey string) (string, error) {
	hashed := createHash(privateKey)

	block, _ := aes.NewCipher([]byte(hashed))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := []byte(hashed)[:12]

	ciphertext := gcm.Seal(nil, nonce, []byte(data), nil)

	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt is decrypt string with
func Decrypt(data string, privateKey string) (string, error) {
	hashed := createHash(privateKey)
	key := []byte(hashed)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := []byte(hashed)[:12]

	plaintext, err := gcm.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}
