package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"bitbucket.org/paydoh/paydoh-commons/httpservice"

	"bankapi/constants"
	"bankapi/requests"
	"bankapi/responses"
)

type NotificationService struct {
	service *httpservice.HttpService
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		service: httpservice.NewHttpService(constants.NotificationURL),
	}
}

func (n *NotificationService) SendNotification(request *requests.NotificationRequest) (*responses.NotificationResponse, error) {

	jsonData, err := json.MarshalIndent(request, "", "   ")

	if err != nil {
		fmt.Println("ERR ", err)
		return nil, err
	}

	response, err := n.service.Post("/api/notify/send", jsonData, map[string]string{
		"Content-Type": "application/json",
	})

	if err != nil {
		fmt.Println("ER ", err)
		return nil, err
	}

	defer response.Body.Close()

	fmt.Println("Sstats ", response.StatusCode)

	if response.StatusCode != http.StatusOK {
		responseData, _ := io.ReadAll(response.Body)
		fmt.Println("error string ", string(responseData))
		return nil, errors.New("failed to send notification")
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	notificationResponse := responses.NewNotificationResponse()

	if err := notificationResponse.Decode(body); err != nil {
		return nil, err
	}

	return notificationResponse, nil
}

func (noti *NotificationService) GetMasterNotification(request string) (res *responses.MasterNotificationRes, er error) {

	response, err := noti.service.Get("/api/notificationMsg/notification/"+request, map[string]string{
		"Content-Type": "application/json",
	})

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to get Master Notification")
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	masterNotificationRes := responses.NewMasterNotificationResponse()

	if err := masterNotificationRes.Decode(body); err != nil {
		return nil, err
	}

	return masterNotificationRes, nil
}

func (noti *NotificationService) SendVerificationEmail(request interface{}) (interface{}, error) {

	jsonData, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	response, err := noti.service.Post("/api/email/send-verification-mail", jsonData, map[string]string{
		"Content-Type": "application/json",
	})

	if err != nil {
		fmt.Println("ER ", err)
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseData, _ := io.ReadAll(response.Body)
		fmt.Println("error string ", string(responseData))
		return nil, errors.New("failed to send notification")
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	notificationResponse := responses.NewNotificationResponse()

	if err := notificationResponse.Decode(body); err != nil {
		return nil, err
	}

	return notificationResponse, nil
}

func (noti *NotificationService) SendAccountInformation(request interface{}) (interface{}, error) {

	jsonData, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	response, err := noti.service.Post("/api/email/send-account-info", jsonData, map[string]string{
		"Content-Type": "application/json",
	})

	if err != nil {
		fmt.Println("ER ", err)
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to send notification")
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	notificationResponse := responses.NewNotificationResponse()

	if err := notificationResponse.Decode(body); err != nil {
		return nil, err
	}

	return notificationResponse, nil
}
