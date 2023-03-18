package utils

import (
	"io"
	"os"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	_ "github.com/joho/godotenv/autoload"
)

func EncryptAES(plaintext []byte) (string, error) {
	key := []byte(os.Getenv("XYZ_SECRET_KEY")) // key must be 32 bit

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	resByte := gcm.Seal(nonce, nonce, plaintext, nil)
	resHex := hex.EncodeToString(resByte)

	return resHex, nil
}

func DecryptAES(str string) ([]byte, error) {
	strKey := os.Getenv("XYZ_SECRET_KEY")
	key := []byte(strKey) // key must be 32 bit

	ciphertext, err := hex.DecodeString(str)
	if err != nil {
		return []byte{}, err
	}
	
	c, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return []byte{}, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}

	return plaintext, nil
}
