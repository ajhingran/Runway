package cheapflight

import (
	"fmt"
	"math"
	"os"
	"time"
)

func ProcessUserRequest() {
	cheapestArgs, excludedAirline, target, err := ProcessArgs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	minFound := math.Inf(1)
	for time.Now().Before(cheapestArgs.RangeStartDate) {
		var message Message
		if cheapestArgs.TripLength == -1 {
			message = GetCheapestOffersFixedDates(cheapestArgs, excludedAirline)
		} else {
			message = GetCheapestOffersRange(cheapestArgs, excludedAirline)
		}

		if message == (Message{}) {
			fmt.Println(fmt.Errorf("unable to find flights at this time"))
		} else {
			if float64(message.Price) < minFound {
				messageString := FormatMessageBody(message)
				SendSMS(messageString)
			} else if float64(message.Price) < target {
				messageString := FormatMessageBodyTarget(message, target)
				SendSMS(messageString)
			}
		}
		time.Sleep(12 * time.Hour)
	}
}
