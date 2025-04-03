package responses

type StaticParameters struct {
	DebitCardAmount  string `json:"debitcard_amt"`
	TollFreeNumber   string `json:"toll_free_number"`
	SupportMailID    string `json:"support_mail_id"`
	AwsCloudFrontUrl string `json:"aws_cloudfront_url"`
}

func NewStaticParameters() *StaticParameters {
	return &StaticParameters{}
}
