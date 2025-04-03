package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keySize       = 32 // 256 bits for AES
	pwdIterations = 65536
	saltSize      = 20            // 20 bytes salt
	blockSize     = aes.BlockSize // 16 bytes for AES
)

type AESEncryptionUtil struct {
	salt []byte
}

func NewAESEncryptionUtil() *AESEncryptionUtil {
	// Use a static salt for consistency with Java implementation
	return &AESEncryptionUtil{
		salt: make([]byte, saltSize),
	}
}

// Apply PKCS5 padding (same as PKCS7 for blockSize 16)
func pkcs5Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padText...)
}

// Remove PKCS5 padding
func pkcs5UnPadding(paddedText []byte) ([]byte, error) {
	length := len(paddedText)
	if length == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(paddedText[length-1])
	if padding > blockSize || padding > length {
		return nil, errors.New("invalid padding")
	}
	return paddedText[:(length - padding)], nil
}

func (a *AESEncryptionUtil) encrypt(plainText, panKey string) (string, error) {
	key := pbkdf2.Key([]byte(panKey), a.salt, pwdIterations, keySize, sha1.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainTextBytes := pkcs5Padding([]byte(plainText), blockSize)

	cipherText := make([]byte, len(plainTextBytes))
	for start := 0; start < len(plainTextBytes); start += blockSize {
		block.Encrypt(cipherText[start:start+blockSize], plainTextBytes[start:start+blockSize])
	}

	// Encode to Base64
	base64Encoded := base64.StdEncoding.EncodeToString(cipherText)
	// Convert Base64 encoded string to Hexadecimal
	hexEncoded := hex.EncodeToString([]byte(base64Encoded))
	return hexEncoded, nil
}

func (a *AESEncryptionUtil) CardDecryption(encryptedText, panKey string) (string, error) {
	key := pbkdf2.Key([]byte(panKey), a.salt, pwdIterations, keySize, sha1.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Decode from Hexadecimal
	encryptedBytes, err := hex.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	// Decode Base64
	base64DecodedBytes, err := base64.StdEncoding.DecodeString(string(encryptedBytes))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(base64DecodedBytes))
	for start := 0; start < len(base64DecodedBytes); start += blockSize {
		block.Decrypt(cipherText[start:start+blockSize], base64DecodedBytes[start:start+blockSize])
	}

	// Unpad and return the decrypted text
	decryptedText, err := pkcs5UnPadding(cipherText)
	if err != nil {
		return "", err
	}

	return string(decryptedText), nil
}
