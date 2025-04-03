package integration

import (
	"bankapi/app"
	"bankapi/constants"
	"bankapi/models"
	"bankapi/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnboardingUserStatus(t *testing.T) {
	router, cleanUpFunc := app.SetupTestRouter()
	defer cleanUpFunc()

	assert := assert.New(t)

	// add onboarding data
	err := models.CreateUserOnboardingStatus(constants.AUTHORIZATION_STEP, "6c372006-4522-4")
	require.NoError(t, err)

	err = models.UpdateUserOnboardingStatus(constants.PERSONAL_DETAILS_STEP, "6c372006-4522-4")
	require.NoError(t, err)

	err = models.UpdateUserOnboardingStatus(constants.DEMOGRAPHIC_FETCH_STAGE, "6c372006-4522-4")
	require.NoError(t, err)

	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmMzNzIwMDYtNDUyMi00fGwpcXJmRzVGMVUkeVQ3cFBWbTJEam1GfjU4Km17UE5OIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjA0MzE3NjN9.GOon18ecMrmGEFqK_mZgDKmHQmTjTkGiusCHT8Nzlmg"

	req, err := http.NewRequest(http.MethodGet, "/api/onboarding/user-status", nil)
	require.NoError(t, err)

	utils.SetRequestHeaders(req, bearerToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println("w.Body.>>>>>>>>", w.Body.String())

	respData := make(map[string]interface{})

	err = json.Unmarshal(w.Body.Bytes(), &respData)
	require.NoError(t, err)

	statusFloat, ok := respData["status"].(float64)
	if !ok {
		t.Fatalf("unexpected type for status: %T", respData["status"])
	}

	assert.Equal(200, int(statusFloat), "status should be equal to 200")
}
