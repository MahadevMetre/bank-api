package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	mathRand "math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dutchcoders/go-clamd"
)

// GenerateRandomCode Function to generate a random SMS authentication code
func GenerateRandomCode(length int) string {
	// Seed the random number generator
	mathRand.Seed(time.Now().UnixNano())

	// Define the length of the authentication code

	// Generate a random code
	code := ""
	for i := 0; i < length; i++ {
		// ASCII values for digits 0-9 are 48-57
		digit := mathRand.Intn(10) + 48
		code += string(rune(digit))
	}

	return code
}

// GenerateRandomUUID generates a random UUID (version 4) without using external libraries.
func GenerateRandomUUID(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length should be greater than 0")
	}

	uuid := make([]byte, 16)

	// Set the version (bits 12-15) to 0100 for UUID version 4.
	uuid[6] = 0x40 | (uuid[6] & 0xF)

	// Set the variant (bits 16-17) to 10 for RFC 4122.
	uuid[8] = 0x80 | (uuid[8] & 0x3F)

	// Fill the first 8 bytes with the current timestamp.
	timestampBytes := uint64ToBytes(uint64(time.Now().UnixNano() / 100))
	copy(uuid[0:8], timestampBytes[2:])

	// Fill the rest of the UUID with random bytes.
	_, err := rand.Read(uuid[8:])
	if err != nil {
		return "", err
	}

	// Convert to UUID string format.
	uuidString := fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])

	// Trim or pad the string to achieve the desired length
	if len(uuidString) >= length {
		return uuidString[:length], nil
	}

	return uuidString + fmt.Sprintf("%0*s", length-len(uuidString), ""), nil
}

func uint64ToBytes(value uint64) []byte {
	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[i] = byte(value >> ((7 - i) * 8) & 0xFF)
	}
	return result
}

// GenerateRandomPassphrase generates a random passphrase.
func GenerateRandomPassphrase() string {
	// Define the character set for the passphrase
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=[]{};:,.<>?/~"

	// Set the length of the passphrase
	passphraseLength := 32

	// Generate the passphrase
	passphrase := make([]byte, passphraseLength)
	charSetLength := big.NewInt(int64(len(charSet)))

	for i := range passphrase {
		randomIndex, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			fmt.Println("Error generating random index:", err)
			return ""
		}
		passphrase[i] = charSet[randomIndex.Int64()]
	}

	return string(passphrase)
}

// Encrypt encrypts the input byte array using AES encryption with the provided key.
//
// input []byte - the input byte array to be encrypted
// key []byte - the key used for encryption
// (string, error) - returns the encrypted string and any error encountered during the encryption process
func Encrypt(input []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, input, nil)

	return hex.EncodeToString(nonce) + hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given ciphertext using the provided key.
//
// It takes a ciphertext string and a key byte array as parameters and returns a string and an error.
func Decrypt(ciphertext string, key []byte) (string, error) {
	decoded, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decoded) < 24 {
		return "", errors.New("invalid ciphertext")
	}

	nonce := decoded[:12]
	ciphertextBytes := decoded[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

const otpChars = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func IsValidAmountFormat(str string) bool {
	regex := `^\d{1,18}(\.\d{2})?$`
	match, _ := regexp.MatchString(regex, str)
	return match
}

func IsValidFileExtension(filename string) (string, bool) {
	allowedExtensions := []string{".jpg", ".jpeg", ".pdf"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range allowedExtensions {
		if ext == validExt {
			return ext, true
		}
	}
	return ext, false
}

func ScanDocument(filePath string) (bool, string, error) {
	c := clamd.NewClamd("tcp://localhost:3310") // ClamAV service address
	f, err := os.Open(filePath)
	if err != nil {
		return false, "", err
	}
	defer f.Close()
	response, err := c.ScanStream(f, make(chan bool))
	if err != nil {
		return false, "", err
	}
	for res := range response {
		if res.Status == clamd.RES_FOUND {
			return false, res.Description, nil
		}
	}
	return true, "", nil
}
