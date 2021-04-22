package twiliogo

import "fmt"

type Error struct {
	Description string
}

func (e Error) Error() string {
	return e.Description
}

type TwilioError struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"more_info"`
}

func (e TwilioError) Error() string {
	var message string

	message = "Twilio Error, "

	if e.Status != 0 {
		message += fmt.Sprintf("Status: %d", e.Status)
	}

	if e.Code != 0 {
		message += fmt.Sprintf(", Code: %d", e.Code)
	}

	if e.Message != "" {
		message += ", Message: " + e.Message
	}

	return message
}
