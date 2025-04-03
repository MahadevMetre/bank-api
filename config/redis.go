package config

import (
	"bankapi/constants"
	"context"
	"log"

	"bitbucket.org/paydoh/paydoh-commons/database"
)

var redis *database.InMemory

func InitRedis(appCtx context.Context) (*database.InMemory, error) {
	memoryConfig := database.InMemoryConfig{
		NetworkType: constants.RedisNetworkType,
		Address:     constants.RedisURL,
		Username:    constants.RedisUserName,
		Password:    constants.RedisPassword,
		DB:          constants.RedisDB,
	}

	memory, err := database.NewInMemory(appCtx, memoryConfig)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	redis = memory

	return memory, nil
}

func GetRedis() *database.InMemory {
	if redis == nil {
		panic("redis is not initialized")
	}
	return redis
}
