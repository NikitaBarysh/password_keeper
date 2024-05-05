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
	client     *http.Client
	hashKey    string
	encrypt    *encryption.Encryptor
	addData    chan []byte
	getData    chan []byte
	deleteData chan []byte
	token      string
	address    string
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
		return nil, err
	}
	return &Sender{
		client:     newClient,
		encrypt:    enc,
		hashKey:    cfg.HashKey,
		addData:    addData,
		getData:    getData,
		deleteData: deleteData,
		address:    cfg.Url,
	}, nil
}

func (s *Sender) PostUserRequest(login, password, path string) error {
	body := entity.User{
		Login:    login,
		Password: password,
	}

	b, err := s.encryptUser(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/%s", s.address, path)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return errors.New("Failed to register user ")
	}

	err = s.parseAuthToken(resp)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) PostDataRequest(data, eventType string) error {
	err := s.checker()
	if err != nil {
		return err
	}

	b, err := encryption.SymmetricEncrypt([]byte(data), s.hashKey)
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("http://%s/api/set/%s", s.address, eventType)

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return errors.New("Status code not 200 ")
	}

	return nil
}

func (s *Sender) GetDataRequest(eventType string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/api/get/%s", s.address, eventType)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, errors.New("Status code not 200 ")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body %w ", err)
	}

	data, err := encryption.SymmetricDecrypt(b, s.hashKey)
	if err != nil {
		return nil, fmt.Errorf("Error decrypting data %w ", err)
	}

	return data, nil
}

func (s *Sender) DeleteDataRequest(eventType string) error {
	url := fmt.Sprintf("http://%s/api/delete/%s", s.address, eventType)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return errors.New("Status code not 200 ")
	}

	return nil
}

func (s *Sender) ConnectWs() error {
	err := s.checker()
	if err != nil {
		return err
	}

	header := http.Header{}
	header.Set("authorization", s.token)

	URL := url.URL{Scheme: "ws", Path: "/ws/connect", Host: "localhost:8080"} // TODO улучшить аддрес

	c, _, err := websocket.DefaultDialer.Dial(URL.String(), header)
	if err != nil {
		return err
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
				break
			}
			go s.getInput()
		case get := <-s.getData:
			err = c.WriteMessage(websocket.TextMessage, get)
			if err != nil {
				break
			}

			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			b, err := encryption.SymmetricDecrypt(msg, s.hashKey)
			if err != nil {
				break
			}

			fmt.Println(string(b))
			go s.getInput()

		case del := <-s.deleteData:
			err = c.WriteMessage(websocket.TextMessage, del)
			if err != nil {
				break
			}
			go s.getInput()
		case <-interrupt:
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return err
			}
			return nil
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
			break
		}
		s.addData <- res
	case "get":
		eventType := input("event")
		str := entity.WebDataMSG{Action: "get", EventType: eventType}
		res, err := s.pack(str)
		if err != nil {
			break
		}
		s.getData <- res
	case "delete":
		eventType := input("event")
		str := entity.WebDataMSG{Action: "delete", EventType: eventType}
		res, err := s.pack(str)
		if err != nil {
			break
		}
		s.deleteData <- res
	}
}

func (s *Sender) pack(data entity.WebDataMSG) ([]byte, error) {
	encData, err := encryption.SymmetricEncrypt(data.Payload, s.hashKey)
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
		return nil, fmt.Errorf("Failed to marshal user to JSON ")
	}

	data, err := s.encrypt.Encrypt(b)
	if err != nil {
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
