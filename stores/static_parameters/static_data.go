package staticparameters

import (
	"bankapi/constants"
	"bankapi/responses"
	"bankapi/utils"
	"context"
	"fmt"
	"time"

	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type Store struct {
	LoggerService *commonSrv.LoggerService
}

func NewStore(log *commonSrv.LoggerService) *Store {

	return &Store{
		LoggerService: log,
	}
}

func (s *Store) GetStaticParameters(ctx context.Context) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/api/static-parameters",
		Message:       "GetStaticParamter log",
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
	}

	res := responses.NewStaticParameters()
	res.DebitCardAmount = constants.DebitCardPaymentAmount
	res.SupportMailID = constants.SupportMailID
	res.TollFreeNumber = constants.TollFreeNumber
	res.AwsCloudFrontUrl = constants.AWSCloudFrontURL

	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = fmt.Sprintf("", res)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return res, nil
}

func (s *Store) GetSecrets(ctx context.Context) (interface{}, error) {
	secrets := map[string]string{
		"razor_pay_id":           constants.RazorPayID,
		"merchant_id":            constants.MERCHANT_ID,
		"long_code_encrypt_text": constants.LONGCODE_ENCRYPT_TEXT,
		"keystore_path_password": constants.KEystorePath,
		"key_alias":              constants.KeyAlias,
		"key_password":           constants.KeyPassword,
		"key_shared_key":         constants.KeySharedKey,
	}

	return secrets, nil
}
