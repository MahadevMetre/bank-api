package models

import (
	"time"

	"bankapi/security"
)

type MessageData struct {
	Id        int64     `json:"id"`
	UserId    string    `json:"user_id"`
	Message   string    `json:"message_data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewMessageData() *MessageData {
	return &MessageData{}
}

func (m *MessageData) Create(userId, message, key string) error {
	m.UserId = userId
	encryptedMessage, err := security.Encrypt([]byte(message), []byte(key))

	if err != nil {
		return err
	}

	m.Message = string(encryptedMessage)

	return nil
}
