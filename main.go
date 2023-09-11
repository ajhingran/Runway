package main

import (
	"fmt"

	runway "github.com/ajhingran/runway/cheapflight"
	sms "github.com/ajhingran/runway/messaging"
)

func main() {
	cheapestArgs, excludedAirline, err := runway.ProcessArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var message runway.Message
	if cheapestArgs.TripLength == -1 {
		message = runway.GetCheapestOffersFixedDates(cheapestArgs, excludedAirline)
	} else {
		message = runway.GetCheapestOffersRange(cheapestArgs, excludedAirline)
	}

	if message == (runway.Message{}) {
		return
	} else {
		messageString := sms.FormatMessageBody(message)
		sms.SendSMS(messageString)
	}
}
