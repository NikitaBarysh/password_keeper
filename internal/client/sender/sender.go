package sender

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"password_keeper/config/client"
	"password_keeper/internal/common/encryption"
	"password_keeper/internal/common/entity"
)

var in = bufio.NewReader(os.Stdin)

const timeOut = 300 * time.Second

type Sender struct {
	//logging       *logger.Logger
	client     *http.Client
	cfg        *client.ClientConfig
	encrypt    *encryption.Encryptor
	addData    chan []byte
	getData    chan []byte
	deleteData chan []byte
	token      string
	//sendInterface SendInterface
}

func NewSender(cfg *client.ClientConfig) (*Sender, error) {
	newClient := &http.Client{
		Timeout: timeOut,
	}

	addData := make(chan []byte, 1)
	getData := make(chan []byte, 1)
	deleteData := make(chan []byte, 1)

	enc, err := encryption.InitEncryptor(cfg.PublicKeyPath)
	if err != nil {
		//logging.Error("Error initializing encryption")
		return nil, err
	}
	return &Sender{
		//logging:       logging,
		client:     newClient,
		encrypt:    enc,
		cfg:        cfg,
		addData:    addData,
		getData:    getData,
		deleteData: deleteData,
		//sendInterface: sendInterface,
	}, nil
}

func (s *Sender) PostUserRequest(login, password, path string) ([]byte, error) {
	body := entity.User{
		Login:    login,
		Password: password,
	}

	b, err := s.encryptUser(body)
	if err != nil {
		//s.logging.Error("Error encrypting user")
		return nil, err
	}

	url := fmt.Sprintf("http://localhost:8080/%s", path) // 	TODO улучшить url

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		//s.logging.Error("Failed to register user")
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		//s.logging.Error("Failed to register user")
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		//s.logging.Error("StatusCode not 200")
		return nil, errors.New("Failed to register user ")
	}

	err = s.parseAuthToken(resp)
	if err != nil {
		//s.logging.Error("Failed to parse auth token")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		//s.logging.Error("Failed to read response body")
		return nil, errors.New("Failed to register user ")
	}

	return respBody, nil
}

func (s *Sender) PostDataRequest(data, eventType string) error {
	err := s.checker()
	if err != nil {
		//s.logging.Error(err.Error())
		return err
	}

	b, err := encryption.SymmetricEncrypt([]byte(data), s.cfg.HashKey)
	if err != nil {
		//s.logging.Error("Error encrypting data")
		return err
	}

	URL := fmt.Sprintf("http://localhost:8080/api/set/%s", eventType)

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(b))
	if err != nil {
		//s.logging.Error("Failed to make request to POST")
		return err
	}

	req.Header.Set("Authorization", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		//s.logging.Error("Failed to add data")
		return err
	}

	if resp.StatusCode/100 != 2 {
		//s.logging.Error("StatusCode not 200")
		return errors.New("Status code not 200 ")
	}

	return nil
}

func (s *Sender) GetDataRequest(eventType string) ([]byte, error) {
	url := fmt.Sprintf("http://localhost:8080/api/get/%s", eventType)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		//s.logging.Error("Failed to make request to GET")
		return nil, err
	}

	req.Header.Set("Authorization", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		//s.logging.Error("Failed to get data")
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		//s.logging.Error("StatusCode not 200")
		return nil, errors.New("Status code not 200 ")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		//s.logging.Error("Failed to read response body")
		return nil, fmt.Errorf("Failed to read response body %w ", err)
	}

	data, err := encryption.SymmetricDecrypt(b, s.cfg.HashKey)
	if err != nil {
		//s.logging.Error("Error decrypting data")
		return nil, fmt.Errorf("Error decrypting data %w ", err)
	}

	return data, nil
}

func (s *Sender) DeleteDataRequest(eventType string) error {
	url := fmt.Sprintf("http://localhost:8080/api/delete/%s", eventType)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		//s.logging.Error("Failed to make request to DELETE")
		return err
	}

	req.Header.Set("Authorization", s.token)

	_, err = s.client.Do(req)
	if err != nil {
		//s.logging.Error("Failed to delete data")
		return err
	}

	//if resp.StatusCode/100 != 2 {
	//	//s.logging.Error("StatusCode not 200")
	//	return errors.New("Status code not 200 ")
	//}

	return nil
}

func (s *Sender) ConnectWs() {
	err := s.checker()
	if err != nil {
		//s.logging.Error("Empty auth token")
	}

	header := http.Header{}
	header.Set("authorization", s.token)

	URL := url.URL{Scheme: "ws", Path: "/ws/connect", Host: "localhost:8080"} // TODO улучшить аддрес

	c, _, err := websocket.DefaultDialer.Dial(URL.String(), header)
	if err != nil {
		//s.logging.Error(err.Error())
		return
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go s.getInput()

	for {
		select {
		case inData := <-s.addData:
			err = c.WriteMessage(websocket.TextMessage, inData)
			if err != nil {
				//s.logging.Error("Failed to write payload")
			}
			go s.getInput()
		case get := <-s.getData:
			err = c.WriteMessage(websocket.TextMessage, get)
			if err != nil {
				//s.logging.Error("Failed to write payload")
			}

			_, msg, err := c.ReadMessage()
			if err != nil {
				//s.logging.Error("Failed to read payload")
			}

			b, err := encryption.SymmetricDecrypt(msg, s.cfg.HashKey)
			if err != nil {
				//s.logging.Error("Error decrypt data")
			}

			fmt.Println(string(b))
			go s.getInput()

		case del := <-s.deleteData:
			err = c.WriteMessage(websocket.TextMessage, del)
			if err != nil {
				//s.logging.Error("Failed to delete ")
			}
			go s.getInput()
		case <-interrupt:
			//s.logging.Info("Caught interrupt signal - quitting!")
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				//s.logging.Error("Failed to close connection")
				return
			}
			return
		}
	}
}

func (s *Sender) getInput() {
	action := input("action")

	switch action {
	case "add":
		data := input("data")
		eventType := input("event")
		str := entity.WebDataMSG{Action: "add", Payload: []byte(data), EventType: eventType}
		res, err := s.pack(str)
		if err != nil {
			//s.logging.Error("Failed to pack data for adding") //  TODO улучшить обработку ошибок
		}
		s.addData <- res
	case "get":
		eventType := input("event")
		str := entity.WebDataMSG{Action: "get", EventType: eventType}
		res, err := s.pack(str)
		if err != nil {
			//s.logging.Error("Failed to pack data for getting")
		}
		s.getData <- res
	case "delete":
		eventType := input("event")
		str := entity.WebDataMSG{Action: "delete", EventType: eventType}
		res, err := s.pack(str)
		if err != nil {
			//s.logging.Error("Failed to pack data for delete")
		}
		s.deleteData <- res
	}
}

func (s *Sender) pack(data entity.WebDataMSG) ([]byte, error) {
	encData, err := encryption.SymmetricEncrypt(data.Payload, s.cfg.HashKey)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting data ")
	}

	data.Payload = encData

	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling data ")
	}

	return b, nil
}

func (s *Sender) encryptUser(user entity.User) ([]byte, error) {
	b, err := json.Marshal(user)
	if err != nil {
		//s.logging.Error("Failed to marshal user to JSON")
		return nil, fmt.Errorf("Failed to marshal user to JSON ")
	}

	data, err := s.encrypt.Encrypt(b)
	if err != nil {
		//s.logging.Error("Failed to encrypt user")
		return nil, fmt.Errorf("Failed to encrypt user ")
	}

	return data, nil
}

func (s *Sender) parseAuthToken(resp *http.Response) error {
	h := resp.Header.Get("Authorization")
	if h == "" {
		return errors.New("Authorization header is empty ")
	}
	s.token = h

	return nil
}

func (s *Sender) checker() error {
	if s.token == "" {
		return errors.New("Token is empty, try to login ")
	}

	return nil
}

func input(cmd string) string {
	fmt.Printf("Write a %s: ", cmd)
	data, err := in.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	data = strings.Replace(data, "\n", "", -1)
	return data
}
