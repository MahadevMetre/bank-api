package integration

import (
	"bankapi/models"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestAccountModel(t *testing.T) {
	db, cleanup := InitTestDatabase(t)
	defer cleanup()
	defer db.Close()

	assert := assert.New(t)

	t.Run("insert account data", func(t *testing.T) {
		// create user personal data
		personalData := &models.PersonalInformation{
			UserId:    "user1",
			Email:     gofakeit.Email(),
			FirstName: gofakeit.Name(),
		}
		if err := models.InsertPersonalInformation(db, personalData); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// insert account data
		if err := models.InsertNewAccount(db, "user1", "test", "test", nil, false, "", "", "", "", ""); err != nil {
			t.Fatalf("insert account data error: %v", err)
		}

		// update account data
		updateData := &models.AccountDataUpdate{
			AccountNumber: "test-account-number",
			CustomerId:    "test-customer_id",
			Status:        "Success",
		}
		if err := models.UpdateAccount(db, updateData, "test"); err != nil {
			t.Fatalf("update account data error: %v", err)
		}

		// get account data
		accountData, err := models.GetAccountDataByUserId(db, "user1")
		if err != nil {
			t.Fatalf("get account data error: %v", err)
		}
		assert.NotNil(accountData, "account data should not be nil")
	})

	t.Run("test get account details", func(t *testing.T) {
		data, err := models.GetAccountDetails(db, "user1")
		if err != nil {
			t.Errorf("account data availability error: %v", err)
		}
		assert.NotNil(data, "account details data should not be nil")
	})

	t.Run("test get account details v2", func(t *testing.T) {
		data, err := models.GetAccountDetailsV2(db, "user1")
		if err != nil {
			t.Errorf("account data availability error: %v", err)
		}
		assert.NotNil(data, "account details data should not be nil")
	})
}
