package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthValues struct {
	UserId    string `json:"user_id"`
	Key       string `json:"key"`
	DeviceIp  string `json:"device_ip"`
	OS        string `json:"os"`
	OSVersion string `json:"os_version"`
	LatLong   string `json:"lat_long"`
}

type RequestPayload struct {
	Payload string `json:"payload"`
}

type TokenData struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Token     string             `bson:"token" json:"token"`
	Expiry    time.Time          `bson:"expiry" json:"expiry"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewTokenData() *TokenData {
	return &TokenData{}
}
