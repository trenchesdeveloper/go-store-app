package notification

type Notification interface {
	SendSMS(phoneNumber, message string) error
}
