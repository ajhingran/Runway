package main

import (
	"encoding/json"
	"fmt"
	runway "github.com/ajhingran/runway/cheapflight"
	"net"
	"os"
	"reflect"
	"strconv"
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

func requestHandler(conn net.Conn) {
	defer conn.Close()
	// Read incoming data
	var user_request UserRequest
	size_buf := make([]byte, 8)

	_, _ = conn.Read(size_buf)
	size_of_req, _ := strconv.Atoi(string(size_buf))
	fmt.Printf("%d\n\n", size_of_req)
	request_buf := make([]byte, size_of_req)
	_, _ = conn.Read(request_buf)

	_ = json.Unmarshal(request_buf, &user_request)

	ureq := reflect.ValueOf(user_request)
	args := []string{os.Args[0]}
	for i := 0; i < ureq.NumField(); i++ {
		field := ureq.Field(i)
		if field.Kind() == reflect.String {
			args = append(args, field.String())
		}
	}
	// Print the incoming data
	go runway.ProcessUserRequest()
}
