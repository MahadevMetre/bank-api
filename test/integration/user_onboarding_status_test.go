package integration

import (
	"bankapi/config"
	"bankapi/constants"
	"bankapi/models"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserOnboardingStatus(t *testing.T) {
	db, cleanup := InitTestDatabase(t)
	defer cleanup()
	defer db.Close()
	config.DB = db

	assert := assert.New(t)

	err := models.CreateUserOnboardingStatus(constants.AUTHORIZATION_STEP, "user1")
	require.NoError(t, err)

	err = models.UpdateUserOnboardingStatus("AUTHORIZATION", "user1")
	require.NoError(t, err)

	err = models.UpdateUserOnboardingStatus("PERSONAL_DETAILS", "user1")
	require.NoError(t, err)

	err = models.UpdateUserOnboardingStatus("DEMOGRAPHIC_FETCH", "user1")
	require.NoError(t, err)

	err = models.UpdateUserOnboardingStatus("ACCOUNT_CREATION", "user1")
	require.NoError(t, err)

	data, err := models.GetUserOnboardingStatus("user1")
	require.NoError(t, err)

	bytesData, err := json.Marshal(data)
	require.NoError(t, err)

	assert.True(data.IsSimVerificationComplete)
	assert.True(data.IsDemographicFetchComplete)
	assert.True(data.IsAccountCreationComplete)

	fmt.Println("data:--------------------------------", string(bytesData))

}
