package main

import (
	"fmt"
	runway "github.com/ajhingran/runway/cheapflight"
	sms "github.com/ajhingran/runway/messaging"
	"math"
	"os"
	"time"
)

func main() {
	cheapestArgs, excludedAirline, target, err := runway.ProcessArgs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	minFound := math.Inf(1)
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
			if float64(message.Price) < minFound {
				messageString := sms.FormatMessageBody(message)
				sms.SendSMS(messageString)
			} else if float64(message.Price) < target {
				messageString := sms.FormatMessageBodyTarget(message, target)
				sms.SendSMS(messageString)
			}
		}
		time.Sleep(12 * time.Hour)
	}
	os.Exit(0)
}
