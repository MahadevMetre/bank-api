package upi

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"bankapi/models"
	"bankapi/security"
)

func ProcessXmlString(m *database.Document, xmlstring, userId string, key string) (string, error) {
	singleResult, err := m.FindOne("upi_token", "upi_token", bson.M{
		"user_id": userId,
	}, bson.M{})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, mongo.ErrNilDocument) {
			b := []byte(xmlstring)

			encrypted, err := security.Encrypt(b, []byte(key))

			if err != nil {
				return "", errors.New("failed to encrypt token")
			}

			if _, err := m.InsertOne("upi_token", "upi_token", bson.M{
				"user_id":    userId,
				"token":      encrypted,
				"expiry":     time.Now().Add(time.Hour * 24 * 90),
				"created_at": time.Now(),
				"updated_at": time.Now(),
			}, false); err != nil {
				return "", err
			}

			return encrypted, nil
		}
		return "", err
	}

	tokenData := models.NewTokenData()

	if err := singleResult.Decode(tokenData); err != nil {
		return "", err
	}

	if tokenData.Expiry.After(time.Now()) {
		b := []byte(xmlstring)

		encrypted, err := security.Encrypt(b, []byte(key))

		if err != nil {
			return "", errors.New("failed to encrypt token")
		}

		if _, err := m.UpdateOne("upi_token", "upi_token", bson.M{
			"user_id": userId,
		}, bson.M{
			"token":      encrypted,
			"expiry":     time.Now().Add(time.Hour * 24 * 90),
			"updated_at": time.Now(),
		}, false, false); err != nil {
			return "", err
		}

		return encrypted, nil
	}

	return tokenData.Token, nil
}

func (s *Store) GenerateCryptoInfo(ctx context.Context, authValues *models.AuthValues) (string, error) {

	clientData, err := s.getClientDataFromRedis(ctx, authValues.UserId)
	if err != nil || len(clientData) == 0 {
		UserClientID, err := models.FindOneClientIDByUserId(s.db, authValues.UserId)
		if err != nil {
			return "", fmt.Errorf("failed to fetch client_id from database: %w", err)
		}
		clientData = map[string]string{
			"client_id":    UserClientID.ClientId.String,
			"login_ref_id": UserClientID.LoginRefId.String,
		}
	}

	clientID := clientData["client_id"]

	serverID := clientData["server_id"]
	if serverID == "" {
		UserServerID, err := models.FindOneServerIDByUserId(s.db, authValues.UserId)
		if err != nil {
			return "", fmt.Errorf("failed to fetch server_id from database: %w", err)
		}
		serverID = UserServerID.ServerId.String
	}

	loginRefID := clientData["login_ref_id"]
	if loginRefID == "" {
		UserClientID, err := models.FindOneClientIDByUserId(s.db, authValues.UserId)
		if err != nil {
			return "", fmt.Errorf("failed to fetch login_ref_id from database: %w", err)
		}
		loginRefID = UserClientID.LoginRefId.String
	}

	existingUserData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		return "", err
	}

	decryptedDeviceId, err := security.Decrypt(existingUserData.DeviceId, []byte(authValues.Key))
	if err != nil {
		return "", err
	}

	cryptoInfo := fmt.Sprintf("%s~%s~%s~%s~%s~%s",
		decryptedDeviceId,
		clientID,
		serverID,
		strings.ToUpper(authValues.OSVersion),
		strings.ToUpper(authValues.OS),
		loginRefID)

	return cryptoInfo, nil
}

func (s *Store) GenerateClientId(platform string) (string, error) {
	switch platform {
	case "ios":
		return generateIOSClientId()
	case "android":
		return generateAndroidClientId(16)
	default:
		return "", fmt.Errorf("unknown platform: %s", platform)
	}
}

func generateAndroidClientId(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length should be greater than 0")
	}

	digits := make([]byte, length)

	firstDigit, err := rand.Int(rand.Reader, big.NewInt(9))
	if err != nil {
		return "", err
	}
	digits[0] = byte(firstDigit.Int64() + 1)

	for i := 1; i < length; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		digits[i] = byte(digit.Int64())
	}

	for i := range digits {
		digits[i] += '0'
	}

	return string(digits), nil
}

func generateIOSClientId() (string, error) {
	uuid := make([]byte, 16)

	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuidStr := fmt.Sprintf("%x", uuid)

	uuidStr = strings.ToUpper(uuidStr)

	uuidStr = strings.ReplaceAll(uuidStr, "-", "")

	if len(uuidStr) < 16 {
		return "", fmt.Errorf("generated UUID string is too short")
	}

	clientID := uuidStr[:16] // Get the first 16 characters

	return clientID, nil
}
