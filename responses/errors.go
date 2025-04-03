package responses

import "encoding/json"

type BankErrorResponse struct {
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
	ErrorCode1    string `json:"rc"`
	ErrorMessage1 string `json:"desc"`
	UpiError      UpiBaseError `json:"Response"`
}

func (b *BankErrorResponse) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b *BankErrorResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, b)
}

func (b *BankErrorResponse) Error() string {
	return b.ErrorMessage
}

type MobileTeamSuccessResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type MobileTeamSuccessResponseWithoutData struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type MobileTeamErrorResponse struct {
    Error struct {
        Errors struct {
            Body string `json:"body"`
        } `json:"errors"`
    } `json:"error"`
    Message string `json:"message"`
    Status  int    `json:"status"`
}
