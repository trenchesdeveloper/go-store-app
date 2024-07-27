package notification

import (
	"encoding/json"
	"fmt"
	"github.com/trenchesdeveloper/go-store-app/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type NotificationClient struct {
	config config.AppConfig
}

func NewNotificationClient(config config.AppConfig) Notification {
	return &NotificationClient{config: config}
}

// SendSMS Twilio
func (n *NotificationClient) SendSMS(phoneNumber, message string) error {
	accountSid := n.config.TwilioAccountSid
	authToken := n.config.TwilioAuthToken
	fromPhone := n.config.TwilioPhoneNumber

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(fromPhone)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
	return nil
}
