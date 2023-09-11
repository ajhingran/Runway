package main

import (
	"fmt"

	runway "github.com/ajhingran/runway/cheapflight"
	"github.com/krisukox/google-flights-api/flights"
)

func main() {
	cheapestArgs, err := runway.ProcessArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	session, err := flights.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	runway.GetCheapestOffersFixedDates(session, cheapestArgs, "")
}
