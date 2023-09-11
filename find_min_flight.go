package main

import (
	"context"
	"errors"
	"fmt"
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
	startDateArg      = 0
	endDateArg        = 1
	durationArg       = 2
	startArg          = 3
	endArg            = 4
	travelerArg       = 5
	classArg          = 6
	tripTypeArg       = 7
	stopArg           = 8
)

func processArgs() (flights.PriceGraphArgs, error) {
	if len(os.Args) < 5 {
		return flights.PriceGraphArgs{}, errors.New("missing minimum number of args")
	}
	args := os.Args[1:]
	startDate, err := time.Parse(defaultDateFormat, args[startDateArg])
	endDate, err := time.Parse(defaultDateFormat, args[endDateArg])

	if err != nil {
		return flights.PriceGraphArgs{}, errors.New("unable to process date fields | mm-dd-yyyy")
	}

	duration, _ := strconv.Atoi(args[durationArg])
	start := strings.Split(args[startArg], "-")
	end := strings.Split(args[endArg], "-")

	if len(start) == 0 || len(end) == 0 {
		return flights.PriceGraphArgs{}, errors.New("need a start and destination city")
	}

	airports := false
	for _, possibleCity := range start {
		if strings.ToUpper(possibleCity) == possibleCity && len(possibleCity) == 3 {
			airports = true
		} else if airports {
			return flights.PriceGraphArgs{}, errors.New("must be all airports in IATA formatting, ie SFO")
		}
	}
	for _, possibleCity := range end {
		if strings.ToUpper(possibleCity) == possibleCity && len(possibleCity) == 3 {
			airports = true
		} else if airports {
			return flights.PriceGraphArgs{}, errors.New("must be all airports in IATA formatting, ie SFO")
		}
	}

	options := flights.Options{
		Travelers: flights.Travelers{Adults: 1},
		Currency:  currency.USD,
		Stops:     flights.AnyStops,
		Class:     flights.Economy,
		TripType:  flights.RoundTrip,
		Lang:      language.English,
	}

	if args[travelerArg] != "default" {
		passengerNum, err := strconv.Atoi(args[travelerArg])
		if err != nil {
			options.Travelers = flights.Travelers{Adults: passengerNum}
		}
	}

	if args[classArg] != "default" {
		class, err := strconv.ParseInt(args[classArg], 10, 64)
		if err != nil {
			switch class {
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

	if args[tripTypeArg] == "OneWay" {
		options.TripType = flights.OneWay
	}

	if args[stopArg] != "default" {
		stops, err := strconv.ParseInt(args[stopArg], 10, 64)
		if err != nil {
			switch stops {
			case int64(flights.Nonstop):
				options.Stops = flights.Nonstop
			case int64(flights.Stop1):
				options.Stops = flights.Stop1
			case int64(flights.Stop2):
				options.Stops = flights.Stop2
			default:
				options.Stops = flights.AnyStops
			}
		}
	}

	var cheapestArgs flights.PriceGraphArgs

	if airports {
		cheapestArgs = flights.PriceGraphArgs{
			RangeStartDate: startDate,
			RangeEndDate:   endDate,
			TripLength:     duration,
			SrcAirports:    start,
			DstAirports:    end,
			Options:        options,
		}
	} else {
		cheapestArgs = flights.PriceGraphArgs{
			RangeStartDate: startDate,
			RangeEndDate:   endDate,
			TripLength:     duration,
			SrcCities:      start,
			DstCities:      end,
			Options:        options,
		}
	}

	return cheapestArgs, nil
}

func getCheapestOffers(session *flights.Session, args flights.PriceGraphArgs) {
	options := args.Options

	priceGraphOffers, err := session.GetPriceGraph(
		context.Background(),
		args,
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, priceGraphOffer := range priceGraphOffers {
		offers, _, err := session.GetOffers(
			context.Background(),
			flights.Args{
				Date:        priceGraphOffer.StartDate,
				ReturnDate:  priceGraphOffer.ReturnDate,
				SrcCities:   args.SrcCities,
				DstCities:   args.DstCities,
				SrcAirports: args.SrcAirports,
				DstAirports: args.DstAirports,
				Options:     options,
			},
		)
		if err != nil {
			fmt.Println(err.Error())
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
			fmt.Println(err.Error())
		}
		if priceRange == nil {
			fmt.Println(err.Error())
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
				fmt.Println(err.Error())
			}
			fmt.Printf("%s %s\n", bestOffer.StartDate, bestOffer.ReturnDate)
			fmt.Printf("price %d\n", int(bestOffer.Price))
			fmt.Println(url)
		}
	}
}

func main() {
	cheapestArgs, err := processArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	session, err := flights.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	getCheapestOffers(session, cheapestArgs)
}
