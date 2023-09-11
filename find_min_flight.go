package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	start := strings.Split(args[3], "-")
	end := strings.Split(args[4], "-")

	if len(start) == 0 || len(end) == 0 {
		return "", errors.New("need a start and destination city")
	}

	airports := false
	for _, possibleCity := range start {
		if strings.ToUpper(possibleCity) == possibleCity && len(possibleCity) == 3 {
			airports = true
		} else if airports {
			return "", errors.New("must be all airports in IATA formatting, ie SFO")
		}
	}
	for _, possibleCity := range end {
		if strings.ToUpper(possibleCity) == possibleCity && len(possibleCity) == 3 {
			airports = true
		} else if airports {
			return "", errors.New("must be all airports in IATA formatting, ie SFO")
		}
	}

	session, err := flights.New()
	if err != nil {
		return "", err
	}

	options := flights.Options{
		Travelers: flights.Travelers{Adults: 1},
		Currency:  currency.USD,
		Stops:     flights.AnyStops,
		Class:     flights.Economy,
		TripType:  flights.RoundTrip,
		Lang:      language.English,
	}

	if args[6] != "default" {
		passengerNum, err := strconv.Atoi(args[6])
		if err != nil {
			options.Travelers = flights.Travelers{Adults: passengerNum}
		}
	}

	if args[7] != "default" {
		class, err := strconv.Atoi(args[7])
		if err != nil {
			switch int64(class) {
			case int64(flights.PremiumEconomy):
				options.Class = flights.PremiumEconomy
			case int64(flights.Business):
				options.Class = flights.Business
			case int64(flights.First):
				options.Class = flights.First
			default:
				options.Class = flights.Economy
			}
		}
	}

	if args[8] != "default" {
		passengerNum, err := strconv.Atoi(args[6])
		if err != nil {
			options.Travelers = flights.Travelers{Adults: passengerNum}
		}
	}

	if args[9] != "default" {
		passengerNum, err := strconv.Atoi(args[6])
		if err != nil {
			options.Travelers = flights.Travelers{Adults: passengerNum}
		}
	}

	if airports {
		getCheapestOffers(session, startDate, endDate, duration, nil, nil, start, end, language.English)
	} else {
		getCheapestOffers(session, startDate, endDate, duration, start, end, nil, nil, language.English)
	}

	return "", err
}

func getCheapestOffers(
	session *flights.Session,
	rangeStartDate, rangeEndDate time.Time,
	tripLength int,
	srcCities, dstCities, srcAirports, dstAirports []string,
	lang language.Tag,
) {
	logger := log.New(os.Stdout, "", 0)

	options := flights.Options{
		Travelers: flights.Travelers{Adults: 1},
		Currency:  currency.USD,
		Stops:     flights.AnyStops,
		Class:     flights.Economy,
		TripType:  flights.RoundTrip,
		Lang:      lang,
	}

	priceGraphOffers, err := session.GetPriceGraph(
		context.Background(),
		flights.PriceGraphArgs{
			RangeStartDate: rangeStartDate,
			RangeEndDate:   rangeEndDate,
			TripLength:     tripLength,
			SrcAirports:    srcAirports,
			DstAirports:    dstAirports,
			Options:        options,
		},
	)
	if err != nil {
		logger.Fatal(err)
	}

	for _, priceGraphOffer := range priceGraphOffers {
		offers, _, err := session.GetOffers(
			context.Background(),
			flights.Args{
				Date:        priceGraphOffer.StartDate,
				ReturnDate:  priceGraphOffer.ReturnDate,
				SrcAirports: srcAirports,
				DstAirports: dstAirports,
				Options:     options,
			},
		)
		if err != nil {
			logger.Fatal(err)
		}

		var bestOffer flights.FullOffer
		for _, o := range offers {
			if o.Price != 0 && (bestOffer.Price == 0 || o.Price < bestOffer.Price) {
				bestOffer = o
			}
		}

		_, priceRange, err := session.GetOffers(
			context.Background(),
			flights.Args{
				Date:        bestOffer.StartDate,
				ReturnDate:  bestOffer.ReturnDate,
				SrcAirports: []string{bestOffer.SrcAirportCode},
				DstAirports: []string{bestOffer.DstAirportCode},
				Options:     options,
			},
		)
		if err != nil {
			logger.Fatal(err)
		}
		if priceRange == nil {
			logger.Fatal("missing priceRange")
		}

		if bestOffer.Price < priceRange.Low {
			url, err := session.SerializeURL(
				context.Background(),
				flights.Args{
					Date:        bestOffer.StartDate,
					ReturnDate:  bestOffer.ReturnDate,
					SrcAirports: []string{bestOffer.SrcAirportCode},
					DstAirports: []string{bestOffer.DstAirportCode},
					Options:     options,
				},
			)
			if err != nil {
				logger.Fatal(err)
			}
			logger.Printf("%s %s\n", bestOffer.StartDate, bestOffer.ReturnDate)
			logger.Printf("price %d\n", int(bestOffer.Price))
			logger.Println(url)
		}
	}
}

func main() {
	t := time.Now()

	session, err := flights.New()
	if err != nil {
		log.Fatal(err)
	}

	getCheapestOffers(
		session,
		time.Now().AddDate(0, 0, 60),
		time.Now().AddDate(0, 0, 90),
		7,
		nil,
		nil,
		[]string{"MSN"},
		[]string{"DCA"},
		language.English,
	)

	fmt.Println(time.Since(t))
}
