package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/krisukox/google-flights-api/flights"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

const (
	defaultDateFormat = "07-04-1999"
)

func processArgs() (string, error) {
	if len(os.Args) < 6 {
		return "", errors.New("missing minimum number of args")
	}
	args := os.Args[1:]
	startDate, err := time.Parse(defaultDateFormat, args[0])
	endDate, err := time.Parse(defaultDateFormat, args[1])

	if err != nil {
		return "", errors.New("unable to process date fields | mm-dd-yyyy")
	}

	duration, _ := strconv.Atoi(args[2])
	getCheapestOffer(startDate, endDate, duration, args[3], args[4], language.English)
	return "", err
}

func getCheapestOffer(
	rangeStartDate, rangeEndDate time.Time,
	tripLength int,
	srcCity, dstCity string,
	lang language.Tag,
) {
	session, err := flights.New()
	if err != nil {
		log.Fatal(err)
	}

	options := flights.Options{
		Travelers: flights.Travelers{Adults: 1},
		Currency:  currency.PLN,
		Stops:     flights.AnyStops,
		Class:     flights.Economy,
		TripType:  flights.RoundTrip,
		Lang:      lang,
	}

	offers, err := session.GetPriceGraph(
		context.Background(),
		flights.PriceGraphArgs{
			RangeStartDate: rangeStartDate,
			RangeEndDate:   rangeEndDate,
			TripLength:     tripLength,
			SrcCities:      []string{srcCity},
			DstCities:      []string{dstCity},
			Options:        options,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	var bestOffer flights.Offer
	for _, o := range offers {
		if o.Price != 0 && (bestOffer.Price == 0 || o.Price < bestOffer.Price) {
			bestOffer = o
		}
	}

	fmt.Printf("%s %s\n", bestOffer.StartDate, bestOffer.ReturnDate)
	fmt.Printf("price %d\n", int(bestOffer.Price))
	url, err := session.SerializeURL(
		context.Background(),
		flights.Args{
			Date:       bestOffer.StartDate,
			ReturnDate: bestOffer.ReturnDate,
			SrcCities:  []string{srcCity},
			DstCities:  []string{dstCity},
			Options:    options,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(url)
}

func main() {
	getCheapestOffer(
		time.Now().AddDate(0, 0, 60),
		time.Now().AddDate(0, 0, 90),
		2,
		"Warsaw",
		"Athens",
		language.English,
	)
}
