// Package sender - пакет для взаимодействия с сервером
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

// timeOut - через какое время выдаст timeOut запросу от клиента
const timeOut = 300 * time.Second

// Sender - структура в которой хранятся свойства для взаимодействия с сервером
type Sender struct {
	client     *http.Client
	hashKey    string
	encrypt    *encryption.Encryptor
	addData    chan []byte
	getData    chan []byte
	deleteData chan []byte
	token      string
	address    string
}

// NewSender - создаем Sender
func NewSender(cfg *client.ClientConfig) (*Sender, error) {
	newClient := &http.Client{
		Timeout: timeOut,
	}

	addData := make(chan []byte, 1)
	getData := make(chan []byte, 1)
	deleteData := make(chan []byte, 1)

	enc, err := encryption.InitEncryptor(cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("NewSender: %w", err)
	}

	newURL := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return &Sender{
		client:     newClient,
		encrypt:    enc,
		hashKey:    cfg.HashKey,
		addData:    addData,
		getData:    getData,
		deleteData: deleteData,
		address:    newURL,
	}, nil
}

// PostUserRequest - в зависимости от типа действия, регистрируемся и авторизуемся
func (s *Sender) PostUserRequest(login, password, path string) error {
	body := entity.User{
		Login:    login,
		Password: password,
	}

	b, err := s.encryptUser(body)
	if err != nil {
		return fmt.Errorf("PostUserRequest: %w", err)
	}

	url := fmt.Sprintf("http://%s/%s", s.address, path)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("PostUserRequest: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("PostUserRequest: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return errors.New("PostUserRequest: Failed to register user: status not 200 ")
	}

	err = s.parseAuthToken(resp)
	if err != nil {
		return fmt.Errorf("PostUserRequest: %w", err)
	}

	return nil
}

// PostDataRequest - отправляем данные на сервер
func (s *Sender) PostDataRequest(data, eventType string) error {
	err := s.checker()
	if err != nil {
		return fmt.Errorf("PostDataRequest: %w", err)
	}

	b, err := encryption.SymmetricEncrypt([]byte(data), s.hashKey)
	if err != nil {
		return fmt.Errorf("PostDataRequest: %w", err)
	}

	URL := fmt.Sprintf("http://%s/api/set/%s", s.address, eventType)

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("PostDataRequest: %w", err)
	}

	req.Header.Set("Authorization", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("PostDataRequest: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return errors.New("PostDataRequest: status code not 200 ")
	}

	return nil
}

// GetDataRequest - получаем данные с сервера
func (s *Sender) GetDataRequest(eventType string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/api/get/%s", s.address, eventType)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetDataRequest: %w", err)
	}

	req.Header.Set("Authorization", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetDataRequest: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return nil, errors.New("GetDataRequest: status code not 200 ")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetDataRequest: Failed to read response body %w ", err)
	}

	data, err := encryption.SymmetricDecrypt(b, s.hashKey)
	if err != nil {
		return nil, fmt.Errorf("GetDataRequest: %w ", err)
	}

	return data, nil
}

// DeleteDataRequest - удаляем данные с сервера
func (s *Sender) DeleteDataRequest(eventType string) error {
	url := fmt.Sprintf("http://%s/api/delete/%s", s.address, eventType)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("DeleteDataRequest: %w", err)
	}

	req.Header.Set("Authorization", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("DeleteDataRequest: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return errors.New("DeleteDataRequest: status code not 200 ")
	}

	return nil
}

// ConnectWs - websocket соединение с сервером
func (s *Sender) ConnectWs() error {
	err := s.checker()
	if err != nil {
		return fmt.Errorf("ConnectWs: %w", err)
	}

	header := http.Header{}
	header.Set("authorization", s.token)

	URL := url.URL{Scheme: "ws", Path: "/ws/connect", Host: "localhost:8080"}

	c, _, err := websocket.DefaultDialer.Dial(URL.String(), header)
	if err != nil {
		return fmt.Errorf("ConnectWs: %w", err)
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
				return fmt.Errorf("ConnectWs: %w", err)
			}
			return nil
		}
	}
}

// getInput - получаем действия от пользователя, которе нудно сделать
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

// pack - упаковываем и шифруем данные для отправки на сервер
func (s *Sender) pack(data entity.WebDataMSG) ([]byte, error) {
	encData, err := encryption.SymmetricEncrypt(data.Payload, s.hashKey)
	if err != nil {
		return nil, fmt.Errorf("pack: %w", err)
	}

	data.Payload = encData

	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("pack: %w", err)
	}

	return b, nil
}

// encryptUser - шифруем логин и пароль
func (s *Sender) encryptUser(user entity.User) ([]byte, error) {
	b, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("encryptUser: %w", err)
	}

	data, err := s.encrypt.Encrypt(b)
	if err != nil {
		return nil, fmt.Errorf("encryptUser: %w", err)
	}

	return data, nil
}

// parseAuthToken - парсим jwt токен, полученный от сервера
func (s *Sender) parseAuthToken(resp *http.Response) error {
	h := resp.Header.Get("Authorization")
	if h == "" {
		return errors.New("parseAuthToken: Authorization header is empty ")
	}
	s.token = h

	return nil
}

// checker - проверяем на наличие токена у пользщователя
func (s *Sender) checker() error {
	if s.token == "" {
		return errors.New("checker: Token is empty, try to login ")
	}

	return nil
}

func (s *Sender) getTokenForTest(path string) (string, error) {
	body := entity.User{
		Login:    "testUserForDataTest",
		Password: "testUserForDataTest",
	}

	b, err := s.encryptUser(body)
	if err != nil {
		return "", fmt.Errorf("PostUserRequest: %w", err)
	}

	url := fmt.Sprintf("http://%s/%s", "localhost:8000", path)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return "", fmt.Errorf("PostUserRequest: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("PostUserRequest: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return "", errors.New("PostUserRequest: Failed to register user: status not 200 ")
	}

	h := resp.Header.Get("Authorization")
	if h == "" {
		return "", errors.New("parseAuthToken: Authorization header is empty ")
	}

	return h, nil
}

// input - считывает команды с консоли
func input(cmd string) string {
	fmt.Printf("Write a %s: ", cmd)
	data, err := in.ReadString('\n')
	if err != nil {
		log.Printf("input %v\n", err)
	}
	data = strings.Replace(data, "\n", "", -1)
	return data
}
