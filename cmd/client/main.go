package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"password_keeper/config/client"
	"password_keeper/internal/client/sender"
	"password_keeper/internal/common/logger"
)

var (
	Version   string
	BuildTime string

	reader *bufio.Reader
)

const (
	register = "register"
	login    = "login"
	password = "password"
)

func main() {
	logging := logger.InitLogger()
	log.Printf("Version: %s, Time: %s", Version, BuildTime)

	cfg := client.NewClient()

	sen, err := sender.NewSender(cfg)
	if err != nil {
		logging.Error(err.Error())
	}

	reader = bufio.NewReader(os.Stdin)

	for {
		action := read("command")

		switch action {
		case register:
			loginUser := read(login)
			passwordUser := read(password)
			err := sen.PostUserRequest(loginUser, passwordUser, register)
			if err != nil {
				logging.Error(err.Error())
				continue
			}
			fmt.Println("Profile created successfully")
		case login:
			loginUser := read(login)
			passwordUser := read(password)
			err := sen.PostUserRequest(loginUser, passwordUser, login)
			if err != nil {
				logging.Error(err.Error())
				continue
			}
			fmt.Println("Profile logged successfully")
		case "add data":
			payload := read("data")
			eventType := read("event")
			err = sen.PostDataRequest(payload, eventType)
			if err != nil {
				logging.Error(err.Error())
				continue
			}
			fmt.Println("Successfully added data")
		case "get data":
			eventType := read("event")
			resp, err := sen.GetDataRequest(eventType)
			if err != nil {
				logging.Error(err.Error())
				continue
			}
			fmt.Println(string(resp))
		case "delete data":
			eventType := read("event")
			err = sen.DeleteDataRequest(eventType)
			if err != nil {
				logging.Error(err.Error())
				continue
			}
			fmt.Println("Successfully deleted data")
		case "ws":
			sen.ConnectWs()
		case "exit":
			return
		}
	}
}

func read(action string) string {
	fmt.Printf("Enter %s: ", action)
	act, err := reader.ReadString('\n')
	act = strings.Replace(act, "\n", "", -1)
	if err != nil {
		fmt.Println(fmt.Errorf("read: Error reading input: %s \n", err))
	}
	return act
}
