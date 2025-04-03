package integration

import (
	"bankapi/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserClientDataModel(t *testing.T) {
	db, cleanup := InitTestDatabase(t)
	defer cleanup()
	defer db.Close()

	assert := assert.New(t)

	t.Run("test insert client data", func(t *testing.T) {
		clientData, err := models.SaveTransIdAndClientIdByUserId(db, "user1", "test-trans-id", "test-client-id")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assert.NotEqual(clientData, nil)
	})

	t.Run("test update client data", func(t *testing.T) {
		updateData, err := models.SaveServerIdByUserId(db, "user1", "testServer-id", "21234", "21234", "21234")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assert.NotEqual(updateData.LoginRefId, "21234", "it should be equal")
	})

}
