package config

import (
	"bankapi/constants"
	"context"
	"fmt"

	"bitbucket.org/paydoh/paydoh-commons/database"
)

func InitMongoDB(appCtx context.Context) (*database.Document, error) {
	client, err := database.ConnectToMongoDB(appCtx, constants.MongodbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	dbdoc, err := database.NewDocument(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB document: %v", err)
	}

	return dbdoc, nil
}
