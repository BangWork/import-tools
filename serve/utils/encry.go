package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

func CBCEncrypt(plaintext string, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	blockSize := len(key)
	padding := blockSize - len(plaintext)%blockSize
	if padding == 0 {
		padding = blockSize
	}

	plaintext += string(bytes.Repeat([]byte{byte(padding)}, padding))
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		return ""
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func CBCDecrypt(ciphertext string, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	cipherCode, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return ""
	}

	iv := cipherCode[:aes.BlockSize]
	cipherCode = cipherCode[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherCode, cipherCode)

	plaintext := string(cipherCode)
	return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])]
}
