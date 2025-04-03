package config

import (
	"bankapi/constants"

	"bitbucket.org/paydoh/paydoh-commons/amazon"
)

func InitAWS() (*amazon.Aws, error) {
	awsKeyId := constants.AWSAccessKeyID
	awsSecretKey := constants.AWSSecretAccessKey
	awsRegion := constants.AWSRegion

	awsInstance, err := amazon.NewAws(
		awsRegion,
		awsKeyId,
		awsSecretKey,
		constants.AWSBucketName,
		constants.AWSCloudFrontURL,
	)

	if err != nil {
		return nil, err
	}

	return awsInstance, nil
}
