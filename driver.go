package main

import (
	"encoding/json"
	runway "github.com/ajhingran/runway/cheapflight"
	"net/http"
	"os"
	"reflect"
)

const (
	address = ":8080"
)

type UserRequest struct {
	RangeStartDate   string `json:"range-start-date"`
	RangeEndDate     string `json:"range-end-date"`
	TripLength       string `json:"trip-length"`
	Src              string `json:"trip-src"`
	Dst              string `json:"trip-dst"`
	Stops            string `json:"stops"`
	Class            string `json:"class"`
	TripType         string `json:"trip-type"`
	ExcludedAirlines string `json:"excluded-airlines"` // Added field for excluded airlines
	Target           string `json:"target"`
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Read incoming data
	var userRequest UserRequest

	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("402 - Unable to process POST"))
		return
	}

	ureq := reflect.ValueOf(userRequest)
	args := []string{os.Args[0]}
	for i := 0; i < ureq.NumField(); i++ {
		field := ureq.Field(i)
		if field.Kind() == reflect.String {
			args = append(args, field.String())
		}
	}
	// Print the incoming data
	go runway.ProcessUserRequest()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request Configured"))
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Runway"))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		processRequest(w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 - Need to issue POST request"))
}

func main() {
	if len(os.Args) != 1 {
		go runway.ProcessUserRequest()
	}

	http.HandleFunc("/", handleHello)
	http.HandleFunc("/request", handleRequest)
	http.HandleFunc("/requests", handleRequest)
	http.ListenAndServe(address, nil)
}
