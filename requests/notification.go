package requests

type NotificationUser struct {
	UserId      string `json:"user_id"`
	DeviceToken string `json:"device_token"`
	OS          string `json:"os"`
	PackageId   string `json:"package_id"`
}

type NotificationTargetRequest struct {
	UserId                string `json:"user-id" validate:"required"`
	DeviceToken           string `json:"device-token" validate:"required"`
	ChannelType           string `json:"channel-type" validate:"required,oneof=apns fcm"`
	ApplicationIdentifier string `json:"application-identifier" validate:"required"`
}

type NotificationRequest struct {
	Targets      []NotificationTargetRequest `json:"targets"`
	Title        string                      `json:"title,omitempty"`
	TypeId       string                      `json:"type-id" validate:"required"`
	Body         string                      `json:"body,omitempty"`
	ImageUrl     string                      `json:"image-url,omitempty"`
	ClickAction  string                      `json:"click-action,omitempty"`
	Sound        string                      `json:"sound,omitempty" validate:"default_sound"`
	Badge        int                         `json:"badge,omitempty" validate:"default_badge"`
	TTL          int                         `json:"ttl,omitempty" validate:"default_ttl"`
	Priority     string                      `json:"priority" validate:"default_priority"`
	ApnsPushType string                      `json:"apns-push-type,omitempty" validate:"default_apns_push"`
	CustomData   map[string]string           `json:"custom-data,omitempty"`
}

func NewNotificationUser() *NotificationUser {
	return &NotificationUser{}
}

func NewNotificationRequest() *NotificationRequest {
	return &NotificationRequest{}
}

func (request *NotificationRequest) CreateNotificationPayload(
	users []NotificationUser,
	title,
	body,
	typeId,
	clickAction string) error {
	targets := make([]NotificationTargetRequest, 0)

	for _, user := range users {

		var channeltype string

		if user.OS == "android" {
			channeltype = "fcm"
		} else {
			channeltype = "apns"
		}

		target := &NotificationTargetRequest{
			UserId:                user.UserId,
			DeviceToken:           user.DeviceToken,
			ChannelType:           channeltype,
			ApplicationIdentifier: user.PackageId,
		}

		targets = append(targets, *target)
	}

	request.Targets = append(request.Targets, targets...)
	request.Title = title
	request.Body = body
	request.TypeId = typeId
	request.ClickAction = clickAction
	request.Priority = "high"
	request.CustomData = map[string]string{
		"verification-status": clickAction,
	}

	return nil
}
