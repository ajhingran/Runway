package cheapflight

import (
	"fmt"
	"math"
	"os"
	"time"
)

func ProcessUserRequest() {
	fmt.Println(os.Args)
	cheapestArgs, excludedAirline, target, SMSNum, err := ProcessArgs()
	if err != nil {
		fmt.Println(err.Error())
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
				SendSMS(messageString, SMSNum)
			} else if float64(message.Price) < target {
				messageString := FormatMessageBodyTarget(message, target)
				SendSMS(messageString, SMSNum)
			}
		}
		time.Sleep(12 * time.Hour)
	}
}
