package utils

import (
	"bankapi/requests"
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type DebitCardControlUtils struct{}

func NewDebitCardControlUtils() *DebitCardControlUtils {
	return &DebitCardControlUtils{}
}

// Pad the input to be a multiple of the block size
func cardPad(src []byte) []byte {
	blockSize := aes.BlockSize
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Unpad the input
func cardUnpad(src []byte) ([]byte, error) {
	padding := src[len(src)-1]
	if int(padding) > len(src) {
		return nil, errors.New("invalid padding")
	}
	return src[:len(src)-int(padding)], nil
}

// Encrypt encrypts plaintext using AES in ECB mode
func cardControlEncrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedText := cardPad(plaintext)
	ciphertext := make([]byte, len(paddedText))
	for i := 0; i < len(paddedText); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], paddedText[i:i+aes.BlockSize])
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES in ECB mode
func cardControlDecrypt(ciphertext string, key []byte) ([]byte, error) {
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

	return cardUnpad(decrypted)
}

func DebitCardGenerateEncryptedReq(req []byte, txnId string, enckey string) (*requests.EncryptedReq, error) {

	key := []byte(enckey)
	if len(key) != 32 {
		return nil, errors.New("Invalid Key")
	}

	result1, err := cardControlEncrypt([]byte(req), key)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	result["input"] = result1
	result["TransactionId"] = txnId
	jsondata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	encrypted, err := cardControlEncrypt(jsondata, key)
	if err != nil {
		return nil, err
	}
	finalReq := requests.EncryptedReq{}
	data, err := finalReq.Bind(encrypted)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func DebitCardDecryptResponse(respo string, encKey string) ([]byte, error) {
	key := []byte(encKey) // Must be 32 bytes for AES-256
	if len(key) != 32 {
		return nil, errors.New("invalid key")
	}
	if respo == "" {
		return nil, errors.New("Response not found")
	}
	decrypt1, err := cardControlDecrypt(respo, key)
	if err != nil {
		return nil, err
	}
	var responseData1 map[string]string
	err = json.Unmarshal([]byte(decrypt1), &responseData1)
	if err != nil {

		return nil, err
	}
	if responseData1["response"] == "" {
		return nil, errors.New("ResponseKey not found")
	}
	decrypt2, err := cardControlDecrypt(responseData1["response"], key)
	if err != nil {

		return nil, err
	}

	return decrypt2, nil
	// return decrypt1, nil
}
