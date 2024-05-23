package sender

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"password_keeper/config/client"
	"password_keeper/internal/common/entity"
)

func TestSenderPostUser(t *testing.T) {
	type mockBehaviour func(s *MockSendInterface)
	tests := []struct {
		name    string
		mock    mockBehaviour
		user    entity.User
		wantErr error
		address string
		action  string
		pubKey  string
	}{
		{
			name:    "Ok",
			mock:    func(s *MockSendInterface) {},
			wantErr: nil,
			address: "localhost:8000",
			user: entity.User{ // Если проверять, то нужно рандомные значение, т.к. эти будут уже в базе
				Login:    "testSenderPost",
				Password: "testSenderPost",
			},
			action: "register",
			pubKey: "./public.rsa",
		},
		{
			name:    "err to do encrypt",
			mock:    func(s *MockSendInterface) {},
			wantErr: errors.New("err encrypt"),
			user: entity.User{
				Login:    "test",
				Password: "test",
			},

			action: "register",
			pubKey: "./public.rsa",
		},
		{
			name:    "status not 200",
			mock:    func(s *MockSendInterface) {},
			wantErr: errors.New("status not 200"),
			address: "localhost:8000",
			user: entity.User{
				Login:    "testSenderPost",
				Password: "testSenderPost",
			},

			action: "register",
			pubKey: "./public.rsa",
		},
		{
			name:    "err do request",
			mock:    func(s *MockSendInterface) {},
			wantErr: errors.New("err do request"),
			user: entity.User{
				Login:    "testdsjadnsfvdsvs",
				Password: "testvsvsvssvsv",
			},
			action: "register",
			pubKey: "./public.rsa",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := &client.ClientConfig{
				Host:          test.address,
				PublicKeyPath: test.pubKey,
			}

			sender, err := NewSender(cfg)
			require.NoError(t, err)

			s := NewMockSendInterface(c)
			test.mock(s)

			if err := sender.PostUserRequest(test.user.Login, test.user.Password, test.action); (err != nil) != (test.wantErr != nil) {
				t.Errorf("PostUserRequest() error = %v, wantErr %v", err, test.wantErr != nil)
			}
		})
	}
}

func TestSenderPostData(t *testing.T) {
	type mockBehaviour func(s *MockSendInterface)
	tests := []struct {
		name      string
		mock      mockBehaviour
		testData  string
		testEvent string
		wantErr   error
		hashKey   string
		host      string
		port      string
		path      string
	}{
		{
			name:      "Ok",
			mock:      func(s *MockSendInterface) {},
			testData:  "hi",
			testEvent: "testPostData",
			wantErr:   nil,
			hashKey:   "some",
			host:      "localhost",
			port:      "8000",
			path:      "register",
		},
		{
			name:     "No event",
			mock:     func(s *MockSendInterface) {},
			testData: "hi",
			wantErr:  errors.New("Add event "),
			hashKey:  "some",
			host:     "localhost",
			port:     "8000",
			path:     "login",
		},
		{
			name:      "No hash func",
			mock:      func(s *MockSendInterface) {},
			testData:  "hi",
			testEvent: "test",
			wantErr:   errors.New("Empty hash key "),
			host:      "localhost",
			port:      "8000",
			path:      "login",
		},
		{
			name:      "Err to do request",
			mock:      func(s *MockSendInterface) {},
			testData:  "hi",
			testEvent: "test",
			hashKey:   "some",
			wantErr:   errors.New("err do request"),
			path:      "login",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := &client.ClientConfig{
				Host:          test.host,
				Port:          test.port,
				PublicKeyPath: "./public.rsa",
				HashKey:       test.hashKey,
			}

			sender, err := NewSender(cfg)
			require.NoError(t, err)

			s := NewMockSendInterface(c)
			test.mock(s)

			token, err := sender.getTokenForTest(test.path)
			require.NoError(t, err)

			sender.token = token

			if err := sender.PostDataRequest(test.testData, test.testEvent); (err != nil) != (test.wantErr != nil) {
				t.Errorf("PostDataRequest() error = %v, wantErr %v", err, test.wantErr != nil)
			}
		})
	}
}

func TestSenderGetData(t *testing.T) {
	type mockBehaviour func(s *MockSendInterface)
	tests := []struct {
		name      string
		mock      mockBehaviour
		testEvent string
		wantErr   error
		hashKey   string
		host      string
		port      string
	}{
		{
			name:      "Ok",
			mock:      func(s *MockSendInterface) {},
			testEvent: "testPostData",
			wantErr:   nil,
			hashKey:   "some",
			host:      "localhost",
			port:      "8000",
		},
		{
			name:      "Empty token",
			mock:      func(s *MockSendInterface) {},
			testEvent: "test",
			wantErr:   errors.New("Token is empty, try to login "),
			hashKey:   "some",
			host:      "localhost",
			port:      "8000",
		},
		{
			name:      "Err to do request",
			mock:      func(s *MockSendInterface) {},
			testEvent: "test",
			wantErr:   errors.New("err do request"),
			hashKey:   "some",
		},
		{
			name:      "Empty hash key",
			mock:      func(s *MockSendInterface) {},
			testEvent: "test",
			wantErr:   errors.New("Empty hash key "),
			host:      "localhost",
			port:      "8000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := &client.ClientConfig{
				Host:          test.host,
				Port:          test.port,
				PublicKeyPath: "./public.rsa",
				HashKey:       test.hashKey,
			}

			sender, err := NewSender(cfg)
			require.NoError(t, err)

			s := NewMockSendInterface(c)
			test.mock(s)

			token, err := sender.getTokenForTest("login")
			require.NoError(t, err)

			sender.token = token

			if _, err := sender.GetDataRequest(test.testEvent); (err != nil) != (test.wantErr != nil) {
				t.Errorf("GetDataRequest() error = %v, wantErr %v", err, test.wantErr != nil)
			}
		})
	}
}

func TestSenderDeleteData(t *testing.T) {
	type mockBehaviour func(s *MockSendInterface)
	tests := []struct {
		name      string
		token     string
		mock      mockBehaviour
		testEvent string
		wantErr   error
		hashKey   string
		address   string
		host      string
		port      string
	}{
		{
			name:      "Ok",
			mock:      func(s *MockSendInterface) {},
			testEvent: "testPostData",
			wantErr:   nil,
			hashKey:   "some",
			host:      "localhost",
			port:      "8000",
		},
		{
			name:      "err to do request",
			mock:      func(s *MockSendInterface) {},
			testEvent: "testPostData",
			wantErr:   errors.New("err to do request"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := &client.ClientConfig{
				Host:          test.host,
				Port:          test.port,
				PublicKeyPath: "./public.rsa",
				HashKey:       test.hashKey,
			}

			sender, err := NewSender(cfg)
			require.NoError(t, err)

			s := NewMockSendInterface(c)
			test.mock(s)

			token, err := sender.getTokenForTest("login")
			require.NoError(t, err)

			sender.token = token

			if err := sender.DeleteDataRequest(test.testEvent); (err != nil) != (test.wantErr != nil) {
				t.Errorf("DeleteDataRequest() error = %v, wantErr %v", err, test.wantErr != nil)
			}
		})
	}
}
