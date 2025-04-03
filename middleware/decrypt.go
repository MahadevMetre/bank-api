package middleware

import (
	"bankapi/requests"
	"bankapi/security"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

func DecryptMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "POST" {
			encryptedRequest := requests.NewEncryptedRequest()
			if err := encryptedRequest.Validate(ctx); err != nil {
				responses.StatusBadRequest(
					ctx,
					customerror.NewError(err),
					"",
				)
				ctx.Abort()
				return
			}

			signingKey := ctx.MustGet("key").(string)
			decrypted, err := security.Decrypt(encryptedRequest.Data, []byte(signingKey))
			if err != nil {
				responses.StatusBadRequest(
					ctx,
					customerror.NewError(err),
					"",
				)
				ctx.Abort()
				return
			}

			ctx.Set("decrypted", decrypted)
			ctx.Next()
		}

		ctx.Next()
	}

}
