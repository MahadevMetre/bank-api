package integration

import (
	"testing"
)

func TestCreateAccount(t *testing.T) {
	// router := app.SetupTestRouter()

	// bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmMzNzIwMDYtNDUyMi00fGwpcXJmRzVGMVUkeVQ3cFBWbTJEam1GfjU4Km17UE5OIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjA0MzE3NjN9.GOon18ecMrmGEFqK_mZgDKmHQmTjTkGiusCHT8Nzlmg"

	// t.Run("CreateAccount", func(t *testing.T) {
	// 	reqData := map[string]interface{}{
	// 		"annual_turn_over":   "3.2",
	// 		"marital_status":     "Single",
	// 		"country_residence":  "IND",
	// 		"mother_maiden_name": "test",
	// 		"customer_education": "test",
	// 		"nationality":        "INDIAN",
	// 	}
	// 	jsonBody, err := json.Marshal(reqData)
	// 	require.NoError(t, err)

	// 	req, err := http.NewRequest(http.MethodPost, "/api/data-encrypt", bytes.NewBuffer(jsonBody))
	// 	require.NoError(t, err)

	// 	req.Header.Add("Authorization", "Bearer "+bearerToken)
	// 	req.Header.Add("X-Device-Ip", "test")
	// 	req.Header.Add("X-OS", "test")
	// 	req.Header.Add("X-OS-Version", "test")
	// 	req.Header.Add("X-Lat-Long", "test,test")

	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	respData := map[string]interface{}{}

	// 	respBody := w.Body.Bytes()
	// 	err = json.Unmarshal(respBody, &respData)
	// 	require.NoError(t, err)

	// 	fmt.Println("w.Body.String():->>>>>>>>>>", respData)

	// 	requestBody := map[string]interface{}{
	// 		"data": respData["data"],
	// 	}
	// 	accountReq, err := json.Marshal(requestBody)
	// 	require.NoError(t, err)

	// req, err = http.NewRequest(http.MethodPost, "/api/onboarding/create-account", bytes.NewBuffer(accountReq))
	// require.NoError(t, err)

	// w = httptest.NewRecorder()

	// // Add the bearer token to the request header
	// bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmMzNzIwMDYtNDUyMi00fGwpcXJmRzVGMVUkeVQ3cFBWbTJEam1GfjU4Km17UE5OIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjA0MzE3NjN9.GOon18ecMrmGEFqK_mZgDKmHQmTjTkGiusCHT8Nzlmg"
	// req.Header.Add("Authorization", "Bearer "+bearerToken)
	// req.Header.Add("X-Device-Ip", "test")
	// req.Header.Add("X-OS", "test")
	// req.Header.Add("X-OS-Version", "test")
	// req.Header.Add("X-Lat-Long", "test,test")

	// router.ServeHTTP(w, req)

	// fmt.Println("w.Body.String():->>", w.Body.String())
	// })
}
