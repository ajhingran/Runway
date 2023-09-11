package main

import (
	"fmt"

	runway "github.com/ajhingran/runway/cheapflight"
)

func main() {
	cheapestArgs, excludedAirline, err := runway.ProcessArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if cheapestArgs.TripLength == -1 {
		runway.GetCheapestOffersFixedDates(cheapestArgs, excludedAirline)
	} else {
		runway.GetCheapestOffersRange(cheapestArgs, excludedAirline)
	}
}
