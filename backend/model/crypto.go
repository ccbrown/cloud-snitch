package model

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
)

func TokenHash(token []byte) []byte {
	h := sha512.Sum512(token)
	return h[:]
}

func NewToken() []byte {
	token := make([]byte, 20)
	if _, err := rand.Read(token); err != nil {
		panic(err)
	}
	return token
}

func pkcs7Trimmed(buf []byte) []byte {
	if len(buf) == 0 || len(buf)%aes.BlockSize != 0 {
		return nil
	}
	n := int(buf[len(buf)-1])
	if n > aes.BlockSize || n <= 0 || len(buf) < n {
		return nil
	}
	return buf[:len(buf)-n]
}

func DecryptSecret(secret, encryptionKey []byte) []byte {
	if len(secret) < aes.BlockSize*2 || len(secret)%aes.BlockSize != 0 {
		return nil
	}
	iv := secret[:aes.BlockSize]
	buf := secret[aes.BlockSize:]
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	return pkcs7Trimmed(buf)
}

func pkcs7Padded(buf []byte) []byte {
	n := aes.BlockSize - (len(buf) % aes.BlockSize)
	ret := make([]byte, len(buf)+n)
	copy(ret, buf)
	copy(ret[len(buf):], bytes.Repeat([]byte{byte(n)}, n))
	return ret
}

func EncryptSecret(secret, encryptionKey []byte) []byte {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	buf := pkcs7Padded(secret)
	mode.CryptBlocks(buf, buf)
	return append(iv, buf...)
}
