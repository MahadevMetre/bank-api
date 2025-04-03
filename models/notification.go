package models

import (
	"bankapi/config"
	"time"

	"bankapi/requests"
	"bankapi/services"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationPreferences struct {
	Prefrence []Prefrences
	UserID    string    `bson:"userid"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Prefrences struct {
	ID               primitive.ObjectID `bson:"_id"`
	NotificationType string             `bson:"notification_type"`
	IsEnable         bool               `bson:"is_enable"`
}

func NewNotificationPreferences() *NotificationPreferences {
	return &NotificationPreferences{}
}

func (s *NotificationPreferences) Bind1(data []AddNotificationPreferences) error {

	for i := 0; i < len(data); i++ {
		list := Prefrences{}
		list.IsEnable = true
		list.NotificationType = data[i].NotificationType
		list.ID = data[i].ID
		s.Prefrence = append(s.Prefrence, list)
	}

	s.UpdatedAt = time.Now()
	s.CreatedAt = time.Now()
	return nil
}

type AddNotificationPreferences struct {
	ID               primitive.ObjectID `bson:"_id"`
	NotificationType string             `bson:"notification_type"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}

func NewAddNotificationPreferences() *AddNotificationPreferences {
	return &AddNotificationPreferences{}
}

func (s *AddNotificationPreferences) Bind(c *gin.Context) error {
	s.ID = primitive.NewObjectID()
	if err := c.BindJSON(s); err != nil {
		return err
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func GenerateNotification(userId string, event, custom, clickAction string) error {

	response, err := services.NewNotificationService().GetMasterNotification(event)
	if err != nil {
		return err
	}
	if response.Data.Status == "1" {
		title := response.Data.Title
		Message := fmt.Sprintf(response.Data.Message, custom)
		db := config.GetDB()
		device, err := FindOneDeviceByUserID(db, userId)
		if err != nil {
			return err
		}

		err = UniversalNotification(device, title, Message, event, clickAction)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func UniversalNotification(devicedata *DeviceData, title, body, typeID, clickAction string) error {

	users := make([]requests.NotificationUser, 0)

	notificationUser := requests.NewNotificationUser()
	notificationUser.UserId = devicedata.UserId
	notificationUser.DeviceToken = devicedata.DeviceToken.String
	notificationUser.PackageId = devicedata.PackageId
	notificationUser.OS = devicedata.OS.String

	users = append(users, *notificationUser)

	notification := requests.NewNotificationRequest()
	err := notification.CreateNotificationPayload(users, title, body, typeID, clickAction)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	_, err = services.NewNotificationService().SendNotification(notification)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return nil
}
