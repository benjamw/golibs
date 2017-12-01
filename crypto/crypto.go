package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"io"
)

func AddSignature(input []byte, key []byte) (output []byte) {
	sig := makeSig(key, input)
	output = append(sig, input...)
	return
}

func CheckSignature(input []byte, key []byte) (remainingBytes []byte, err error) {
	remainingBytes = input[sha1.Size:]
	sig := makeSig(key, remainingBytes)
	if !hmac.Equal(sig, input[:sha1.Size]) {
		return nil, errors.New("signature mismatch")
	}
	return
}

func makeSig(k, b []byte) []byte {
	mac := hmac.New(sha1.New, k)
	mac.Write(b)
	sig := mac.Sum(nil)
	return sig
}

func Encrypt(input []byte, key []byte) ([]byte, error) {
	cryptText := make([]byte, aes.BlockSize+len(input))
	iv := cryptText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	crypt, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(crypt, iv)
	cfb.XORKeyStream(cryptText[aes.BlockSize:], input)
	return cryptText, nil
}

func Decrypt(input []byte, key []byte) ([]byte, error) {
	// clone the slice because input gets mangled by XORKeyStream otherwise
	clone := make([]byte, len(input))
	copy(clone, input)

	if len(clone) < aes.BlockSize {
		return nil, errors.New("input too short")
	}
	crypt, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := clone[:aes.BlockSize]
	cryptText := clone[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(crypt, iv)
	cfb.XORKeyStream(cryptText, cryptText)
	return cryptText, nil
}
