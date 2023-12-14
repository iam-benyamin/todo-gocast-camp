// package client
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"todo-gocast-camp/delivery/deliveryparam"
)

func main() {
	message := os.Args[1]

	connection, err := net.Dial("tcp", "127.0.0.1:1986")
	if err != nil {
		log.Fatalln("error happen")
	}
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {

		}
	}(connection)
	fmt.Println("connection local address", connection.LocalAddr())

	req := deliveryparam.Request{Command: message}
	if req.Command == "create-task" {
		req.CreateTaskRequest = deliveryparam.CreateTaskRequest{
			Title:    "test",
			Duedate:  "test-date",
			Category: 1,
		}
	}

	serializedData, mErr := json.Marshal(&req)
	if mErr != nil {
		log.Fatalln("can't marshal request")
	}

	numberOfWriteBytes, wErr := connection.Write(serializedData)
	if wErr != nil {
		log.Fatalln("can't write data to connection", wErr)
	}

	fmt.Println("number of write bytes", numberOfWriteBytes)

	var data = make([]byte, 1024)
	numberOfReadBytes, rErr := connection.Read(data)
	if rErr != nil {
		log.Println("can't read data from connection", rErr)
	}
	fmt.Println("numberOfReadBytes", numberOfReadBytes)
	fmt.Println("data is: ", string(data))
}
