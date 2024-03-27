package cheapflight

import (
	"fmt"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"os"
)

func SendSMS(alertMessage string) {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(os.Getenv("TO_PHONE_NUMBER"))
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(alertMessage)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("SMS sent successfully!")
	}
}

func FormatMessageBody(m Message) string {
	message := fmt.Sprintf("Lowest offer found at: price %d USD\n"+
		"Flying out on %s\n"+
		"Returning on %s\n"+
		"Check it out here: %s", m.Price, m.Start, m.End, m.Url)
	fmt.Println(message)
	return message
}

func FormatMessageBodyTarget(m Message, target float64) string {
	message := fmt.Sprintf("Flight under target %.2f: price %d USD\n"+
		"Flying out on %s\n"+
		"Returning on %s\n"+
		"Check it out here: %s", target, m.Price, m.Start, m.End, m.Url)
	fmt.Println(message)
	return message
}
