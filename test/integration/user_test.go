package integration

import (
	"fmt"
	"testing"

	"bankapi/config"
	"bankapi/models"
	"bankapi/security"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserAndAccountDetailByMobileNumber(t *testing.T) {
	db, cleanup := InitTestDatabase(t)
	defer cleanup()
	defer db.Close()

	config.DB = db

	assert := assert.New(t)

	t.Run("GetUserAndAccountDetailByMobileNumber", func(t *testing.T) {
		user := models.NewUserData()

		mobileNumber := gofakeit.Phone()

		userid, err := security.GenerateRandomUUID(15)
		require.NoError(t, err)

		rand, err := security.GenerateRandomUUID(6)
		require.NoError(t, err)

		applicantId := fmt.Sprintf("PAYDOH%s", rand)

		passphrase := security.GenerateRandomPassphrase()
		encryptedPassphrase, err := security.Encrypt([]byte(passphrase), []byte("F7846B274E7AAF344BF17C7B6E7DBTSD"))
		require.NoError(t, err)

		if err := db.QueryRow(
			`
			INSERT into user_data (user_id, mobile_number, applicant_id, signing_key) VALUES($1, $2, $3, $4)
			RETURNING id, user_id, mobile_number, signing_key, created_at, updated_at;
				`,
			userid,
			mobileNumber,
			applicantId,
			encryptedPassphrase,
		).Scan(
			&user.Id,
			&user.UserId,
			&user.MobileNumber,
			&user.SigningKey,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			t.Fatalf("error getting user %v", err)
		}

		// create user personal data
		personalData := &models.PersonalInformation{
			UserId:      userid,
			Email:       gofakeit.Email(),
			FirstName:   gofakeit.Name(),
			DateOfBirth: "",
			Gender:      "male",
		}
		if err := models.InsertPersonalInformation(db, personalData); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// insert account data
		if err := models.InsertNewAccount(db, userid, "12312312", "123", nil, false, "", "", "", "", ""); err != nil {
			t.Fatalf("insert account data error: %v", err)
		}

		updateData := &models.AccountDataUpdate{
			UpiId:         "test@upi.kvb",
			AccountNumber: "12312312",
			CustomerId:    "12312312",
			Status:        "success",
		}
		if err := models.UpdateAccountByUserId(updateData, userid); err != nil {
			t.Fatalf("update account data error: %v", err)
		}

		result, err := models.GetUserAndAccountDetailByMobileNumber(db, mobileNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assert.NotNil(t, result)
	})
}
