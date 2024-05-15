package cheapflight

import (
	"fmt"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"net/smtp"
	"os"
)

const (
	smtpserver = "smtp.gmail.com"
	smtpport   = "587"
	smtppair   = smtpserver + ":" + smtpport
)

func SendSMS(alertMessage string, recipient string) {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(recipient)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(alertMessage)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("SMS sent successfully!")
	}
}

func SendEmail(alertMessage string, recipient []string) {
	sender := os.Getenv("FROM_EMAIL")
	pwd := os.Getenv("PWD")
	auth := smtp.PlainAuth("", sender, pwd, smtpserver)
	err := smtp.SendMail(smtppair, auth, sender, recipient, []byte(alertMessage))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Email sent successfully!")
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
