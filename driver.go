package main

import (
	"encoding/json"
	"fmt"
	runway "github.com/ajhingran/runway/cheapflight"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
)

const (
	protocol = "tcp"
	address  = ":8080"
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
	var userRequest UserRequest
	sizeBuf := make([]byte, 4)

	_, err := conn.Read(sizeBuf)
	if err != nil {
		conn.Write([]byte("Error in reading request size"))
		return
	}

	sizeOfReq, _ := strconv.Atoi(string(sizeBuf))
	fmt.Printf("%d\n\n", sizeOfReq)
	requestBuf := make([]byte, sizeOfReq)
	_, err = conn.Read(requestBuf)
	if err != nil {
		conn.Write([]byte("Error in reading request"))
		return
	}

	err = json.Unmarshal(requestBuf, &userRequest)
	if err != nil {
		conn.Write([]byte("Request format incorrect"))
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
}

func main() {
	if len(os.Args) != 1 {
		go runway.ProcessUserRequest()
	}

	port, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	for {
		conn, err := port.Accept()
		if err != nil {
			continue
		}
		go requestHandler(conn)
	}
}
