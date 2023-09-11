package main

import (
	"fmt"
	"time"

	runway "github.com/ajhingran/runway/cheapflight"
	sms "github.com/ajhingran/runway/messaging"
)

func main() {
	cheapestArgs, excludedAirline, err := runway.ProcessArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for time.Now().Before(cheapestArgs.RangeStartDate) {
		var message runway.Message
		if cheapestArgs.TripLength == -1 {
			message = runway.GetCheapestOffersFixedDates(cheapestArgs, excludedAirline)
		} else {
			message = runway.GetCheapestOffersRange(cheapestArgs, excludedAirline)
		}

		if message == (runway.Message{}) {
			fmt.Println(fmt.Errorf("unable to find flights at this time"))
		} else {
			messageString := sms.FormatMessageBody(message)
			sms.SendSMS(messageString)
		}
		time.Sleep(12 * time.Hour)
	}
}
