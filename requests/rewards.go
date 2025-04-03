package requests

import (
	"encoding/json"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type RewardsRequest struct {
	TranId     string  `json:"tran_id" validate:"required"`
	BenfUserId string  `json:"benf_user_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Remarks    string  `json:"remark" validate:"required"`
}

func NewRewardsRequest() *RewardsRequest {
	return &RewardsRequest{}
}

func (r *RewardsRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, r); err != nil {
		return err
	}

	// if !security.IsValidAmountFormat(r.Amount) {
	// 	return errors.New("acount should be of maximum 18 digits with 2 decimal places")
	// }

	return nil
}

func (r *RewardsRequest) ValidateEncrypted(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}
