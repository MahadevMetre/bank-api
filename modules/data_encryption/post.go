package data_encryption

import (
	"bankapi/constants"
	"bankapi/security"
	"bankapi/stores"
	"bankapi/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

func DataEncryption(c *gin.Context) {
	authValues, err := stores.GetAuthValue(c)
	if err != nil {
		fmt.Println("ATUH ERR ", err)
		responses.StatusForbidden(
			c,
			customerror.NewError(err),
		)
		return
	}

	var reqData map[string]interface{}
	err = c.BindJSON(&reqData)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to bind request body")
		return
	}

	bytes, err := json.Marshal(reqData)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to map request body data")
		return
	}

	encString, err := security.Encrypt(bytes, []byte(authValues.Key))
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to encrypt request body data")
		return
	}

	responses.StatusOk(
		c,
		gin.H{
			"encrypted_string": encString,
		},
		"successfully data encrypted.",
		"",
	)
}

func DataDecryption(c *gin.Context) {
	authValues, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusForbidden(
			c,
			customerror.NewError(err),
		)
		return
	}

	if gin.Mode() == gin.ReleaseMode {
		c.Status(http.StatusNotFound)
		return
	}

	var reqData map[string]interface{}
	err = c.BindJSON(&reqData)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to bind request body")
		return
	}

	value := reqData["data"].(string)
	fmt.Println("value:->", value)
	data, err := security.Decrypt(value, []byte(authValues.Key))
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to decrypt request body data")
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	responses.StatusOk(
		c,
		gin.H{
			"decrypted_data": result,
		},
		"successfully data encrypted.",
		"",
	)
}

func BankEncryption(c *gin.Context) {
	// authValues, err := stores.GetAuthValue(c)
	// if err != nil {
	// 	fmt.Println("ATUH ERR ", err)
	// 	responses.StatusForbidden(
	// 		c,
	// 		customerror.NewError(err),
	// 	)
	// 	return
	// }

	var reqData map[string]interface{}
	err := c.BindJSON(&reqData)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to bind request body")
		return
	}

	bytes, err := json.Marshal(reqData)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to map request body data")
		return
	}

	encString, err := utils.GenerateEncryptedReqV2(bytes, constants.BankEncryptionKey)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to encrypt request body data")
		return
	}

	responses.StatusOk(
		c,
		gin.H{
			"encrypted_string": string(encString),
		},
		"successfully data encrypted.",
		"",
	)
}

func BankDecryption(c *gin.Context) {
	// authValues, err := stores.GetAuthValue(c)
	// if err != nil {
	// 	fmt.Println("ATUH ERR ", err)
	// 	responses.StatusForbidden(
	// 		c,
	// 		customerror.NewError(err),
	// 	)
	// 	return
	// }

	var reqData map[string]interface{}
	err := c.BindJSON(&reqData)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to bind request body")
		return
	}
	fmt.Println("\n\nvalue:->")
	value := reqData["data"].(string)
	fmt.Println("\n\nvalue:->", value)
	data, err := utils.DecryptResponse(value, constants.BankEncryptionKey)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to decrypt request body data")
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	responses.StatusOk(
		c,
		gin.H{
			"decrypted_data": result,
		},
		"successfully data encrypted.",
		"",
	)
}
