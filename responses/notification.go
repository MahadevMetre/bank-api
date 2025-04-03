package responses

import "encoding/json"

type NotificationResponse struct {
	Status  int                      `json:"status"`
	Message string                   `json:"message"`
	Data    NotificationResponseData `json:"data"`
}

type NotificationResponseData struct {
	Statuses       []Statuses `json:"statuses"`
	NotificationID string     `json:"notification-id"`
}

type Statuses struct {
	UserID      string `json:"user-id"`
	DeviceToken string `json:"device-token"`
	MessageID   string `json:"message-id"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
}

func NewNotificationResponse() *NotificationResponse {
	return &NotificationResponse{}
}

func (response *NotificationResponse) Decode(data []byte) error {
	if err := json.Unmarshal(data, response); err != nil {
		return err
	}

	return nil
}

type MasterNotificationRes struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    Data   `json:"data"`
}

type Data struct {
	ID        string `json:"_id"`
	CreatedAt string `json:"createdAt"`
	EventType string `json:"eventtype"`
	Message   string `json:"message"`
	Title     string `json:"title"`
	StartDate string `json:"startDate" `
	EndDate   string `json:"endDate"`
	Frequency string `json:"frequency"`
	Status    string `json:"status"` // 0 Inactive 1 Active
	UpdatedAt string `json:"updatedAt"`
	WordCount int    `json:"wordcount"`
}

func NewMasterNotificationResponse() *MasterNotificationRes {
	return &MasterNotificationRes{}
}

func (response *MasterNotificationRes) Decode(data []byte) error {

	if err := json.Unmarshal(data, response); err != nil {
		return err
	}

	return nil
}
