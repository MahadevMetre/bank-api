package utils

import (
	"bankapi/requests"
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type CardControlUtils struct{}

func NewCardControl() *CardControlUtils {
	return &CardControlUtils{}
}

// Pad the input to be a multiple of the block size
func pad(src []byte) []byte {
	blockSize := aes.BlockSize
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Unpad the input
func unpad(src []byte) ([]byte, error) {
	padding := src[len(src)-1]
	if int(padding) > len(src) {
		return nil, errors.New("invalid padding")
	}
	return src[:len(src)-int(padding)], nil
}

// Encrypt encrypts plaintext using AES in ECB mode
func encrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedText := pad(plaintext)
	ciphertext := make([]byte, len(paddedText))
	for i := 0; i < len(paddedText); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], paddedText[i:i+aes.BlockSize])
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES in ECB mode
func decrypt(ciphertext string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {

		return nil, errors.New("failed to decode base64: %v " + err.Error())
	}
	if len(ciphertextBytes)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	decrypted := make([]byte, len(ciphertextBytes))
	for i := 0; i < len(ciphertextBytes); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], ciphertextBytes[i:i+aes.BlockSize])
	}

	return unpad(decrypted)
}

func GenerateEncryptedReq(req interface{}, enckey string) (*requests.EncryptedReq, error) {

	key := []byte(enckey)
	if len(key) != 32 {
		return nil, errors.New("Invalid Key")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	result1, err := encrypt([]byte(body), key)
	if err != nil {
		return nil, err
	}

	finalReq := requests.EncryptedReq{}
	data, err := finalReq.Bind(result1)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func GenerateEncryptedReqV2(reqData []byte, encKey string) ([]byte, error) {
	key := []byte(encKey)

	if len(key) != 32 {
		return nil, errors.New("invalid key")
	}

	result1, err := encrypt([]byte(reqData), key)
	if err != nil {
		return nil, err
	}

	finalReq := requests.EncryptedReq{}
	data, err := finalReq.Bind(result1)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func DecryptResponse(resp string, encKey string) ([]byte, error) {
	if resp == "" {
		return nil, errors.New("response value can't be empty")
	}

	key := []byte(encKey) // Must be 32 bytes for AES-256
	if len(key) != 32 {
		return nil, errors.New("invalid key")
	}

	decrypt1, err := decrypt(resp, key)
	if err != nil {
		return nil, err
	}

	return decrypt1, nil
}
