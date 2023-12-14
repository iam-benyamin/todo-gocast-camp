// package tcpserver
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"todo-gocast-camp/delivery/deliveryparam"
	"todo-gocast-camp/repository/memorystore"
	"todo-gocast-camp/service/task"
)

func main() {
	const (
		network = "tcp"
		address = ":1986"
	)

	listener, lErr := net.Listen(network, address)
	if lErr != nil {
		log.Fatalln("can't listen on given addresses", address, lErr)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Println("can't close listener")
		}
	}(listener)
	fmt.Printf("listen on http://%s:\n", listener.Addr())

	taskMemoryRepo := memorystore.NewTaskStore()
	taskCategoryRepo := memorystore.TaskCategory{
		Task:     taskMemoryRepo,
		Category: nil,
	}

	taskService := task.NewService(taskCategoryRepo)

	for {
		connection, aErr := listener.Accept()
		if aErr != nil {
			log.Println("can't listen to new connection", aErr)

			continue
		}
		// resource leak
		// defer connection.Close()

		var rawRequest = make([]byte, 1024)
		numberOfReadBytes, rErr := connection.Read(rawRequest)
		if rErr != nil {
			log.Println("can't read data from connection", rErr)

			continue
		}
		fmt.Println("data: ", string(rawRequest))

		req := &deliveryparam.Request{}
		if uErr := json.Unmarshal(rawRequest[:numberOfReadBytes], req); uErr != nil {
			log.Println("bad request", uErr)

			continue
		}

		switch req.Command {
		case "create-task":
			response, cErr := taskService.Create(task.CreateRequest{
				Title:               req.CreateTaskRequest.Title,
				Category:            req.CreateTaskRequest.Category,
				Duedate:             req.CreateTaskRequest.Duedate,
				AuthenticatedUserID: 0,
			})
			if cErr != nil {
				_, wErr := connection.Write([]byte(cErr.Error()))
				if wErr != nil {
					log.Println("can't write data to connection", wErr)

					continue
				}
			}

			data, mErr := json.Marshal(&response)
			if mErr != nil {
				_, wErr := connection.Write([]byte(mErr.Error()))
				if wErr != nil {
					log.Println("can't marshal response", wErr)

					continue
				}
			}

			_, wErr := connection.Write([]byte(data))
			if wErr != nil {
				log.Println("can't marshal response", wErr)

				continue
			}
		}
	}
}
