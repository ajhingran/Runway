package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type req struct {
	RangeStartDate   string `json:"range-start-date"`
	RangeEndDate     string `json:"range-end-date"`
	TripLength       string `json:"trip-length"`
	Src              string `json:"trip-src"`
	Dst              string `json:"trip-dst"`
	Travelers        string `json:"travelers"`
	Class            string `json:"class"`
	TripType         string `json:"trip-type"`
	Stops            string `json:"stops"`
	ExcludedAirlines string `json:"excluded-airlines"` // Added field for excluded airlines
	Target           string `json:"target"`
}

func main() {
	newRequest := req{
		RangeStartDate:   "04-11-2024",
		RangeEndDate:     "04-15-2024",
		TripLength:       "-1",
		Src:              "MSN",
		Dst:              "DCA",
		Travelers:        "default",
		Class:            "default",
		TripType:         "default",
		Stops:            "0",
		ExcludedAirlines: "default",
		Target:           "default",
	}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(&newRequest)
	resp, _ := http.Post("http://localhost:8080/request", "application/json", b)
	bytes, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bytes))
}
