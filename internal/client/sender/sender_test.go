package sender

// Для теста для начала нужно запустить сервер, получить токен авторизованного пользователя и поменять токен на полученный
//const tok = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE3Njc2NTMsIlVzZXJJRCI6M30.L_nb24gCL163FFNJsfgvIsqY3u_dho2f6VdJLySEVe8"
//
//func TestSenderPostData(t *testing.T) {
//	//const tok = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE3Njc2NTMsIlVzZXJJRCI6M30.L_nb24gCL163FFNJsfgvIsqY3u_dho2f6VdJLySEVe8"
//	type mockBehaviour func(s *MockSendInterface)
//	tests := []struct {
//		name      string
//		token     string
//		mock      mockBehaviour
//		testData  string
//		testEvent string
//		wantErr   error
//		hashKey   string
//		address   string
//	}{
//		{
//			name:      "Ok",
//			token:     tok,
//			mock:      func(s *MockSendInterface) {},
//			testData:  "hi",
//			testEvent: "testPostData",
//			wantErr:   nil,
//			hashKey:   "some",
//			address:   "localhost:8000",
//		},
//		{
//			name:      "Empty token",
//			token:     "",
//			mock:      func(s *MockSendInterface) {},
//			testData:  "hi",
//			testEvent: "test",
//			wantErr:   errors.New("Token is empty, try to login "),
//			hashKey:   "some",
//			address:   "localhost:8000",
//		},
//		{
//			name:     "No event",
//			token:    tok,
//			mock:     func(s *MockSendInterface) {},
//			testData: "hi",
//			wantErr:  errors.New("Add event "),
//			hashKey:  "some",
//			address:  "localhost:8000",
//		},
//		{
//			name:      "No hash func",
//			token:     tok,
//			mock:      func(s *MockSendInterface) {},
//			testData:  "hi",
//			testEvent: "test",
//			wantErr:   errors.New("Empty hash key "),
//			address:   "localhost:8000",
//		},
//		{
//			name:      "Err to do request",
//			token:     tok,
//			mock:      func(s *MockSendInterface) {},
//			testData:  "hi",
//			testEvent: "test",
//			hashKey:   "some",
//			wantErr:   errors.New("err do request"),
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			s := NewMockSendInterface(c)
//			test.mock(s)
//
//			sss := &Sender{token: test.token, hashKey: test.hashKey, client: http.DefaultClient,
//				address: test.address}
//
//			if err := sss.PostDataRequest(test.testData, test.testEvent); (err != nil) != (test.wantErr != nil) {
//				t.Errorf("PostDataRequest() error = %v, wantErr %v", err, test.wantErr != nil)
//			}
//		})
//	}
//}
//
//func TestSenderPostUser(t *testing.T) {
//	//const tok = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA5ODI4MzUsIlVzZXJJRCI6MX0.MNJez9Bw0MKpzZFoWWjawpmtbmg9AZ-5qQMOQ7PIYpk"
//	type mockBehaviour func(s *MockSendInterface)
//	tests := []struct {
//		name    string
//		token   string
//		mock    mockBehaviour
//		user    entity.User
//		wantErr error
//		address string
//		action  string
//		pubKey  string
//	}{
//		{
//			name:    "Ok",
//			mock:    func(s *MockSendInterface) {},
//			wantErr: nil,
//			address: "localhost:8000",
//			user: entity.User{ // Если проверять, то нужно рандомные значение, т.к. эти будут уже в базе
//				Login:    strconv.Itoa(rand.Int()),
//				Password: "testvsvsvdcsacasaaacsazxxaasvsv",
//			},
//			action: "register",
//			pubKey: "./public.rsa",
//		},
//		{
//			name:    "err to do encrypt",
//			mock:    func(s *MockSendInterface) {},
//			wantErr: errors.New("err encrypt"),
//			user: entity.User{
//				Login:    "test",
//				Password: "test",
//			},
//
//			action: "register",
//			pubKey: "./public.rsa",
//		},
//		{
//			name:    "status not 200",
//			mock:    func(s *MockSendInterface) {},
//			wantErr: errors.New("status not 200"),
//			address: "localhost:8000",
//			user: entity.User{
//				Login:    "test",
//				Password: "test",
//			},
//
//			action: "register",
//			pubKey: "./public.rsa",
//		},
//		{
//			name:    "err do request",
//			mock:    func(s *MockSendInterface) {},
//			wantErr: errors.New("err do request"),
//			user: entity.User{
//				Login:    "testdsjadnsfvdsvs",
//				Password: "testvsvsvssvsv",
//			},
//			action: "register",
//			pubKey: "./public.rsa",
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			cfg := &client.ClientConfig{
//				Url:           test.address,
//				PublicKeyPath: test.pubKey,
//			}
//
//			sender, err := NewSender(cfg)
//			require.NoError(t, err)
//
//			s := NewMockSendInterface(c)
//			test.mock(s)
//
//			if err := sender.PostUserRequest(test.user.Login, test.user.Password, test.action); (err != nil) != (test.wantErr != nil) {
//				t.Errorf("PostUserRequest() error = %v, wantErr %v", err, test.wantErr != nil)
//			}
//		})
//	}
//}
//
//func TestSenderGetData(t *testing.T) {
//	//const tok = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE3Njc2NTMsIlVzZXJJRCI6M30.L_nb24gCL163FFNJsfgvIsqY3u_dho2f6VdJLySEVe8"
//	type mockBehaviour func(s *MockSendInterface)
//	tests := []struct {
//		name      string
//		token     string
//		mock      mockBehaviour
//		testEvent string
//		wantErr   error
//		hashKey   string
//		address   string
//	}{
//		{
//			name:      "Ok",
//			token:     tok,
//			mock:      func(s *MockSendInterface) {},
//			testEvent: "testPostData",
//			wantErr:   nil,
//			hashKey:   "some",
//			address:   "localhost:8000",
//		},
//		{
//			name:      "Empty token",
//			mock:      func(s *MockSendInterface) {},
//			testEvent: "test",
//			wantErr:   errors.New("Token is empty, try to login "),
//			hashKey:   "some",
//			address:   "localhost:8000",
//		},
//		{
//			name:      "Err to do request",
//			token:     tok,
//			mock:      func(s *MockSendInterface) {},
//			testEvent: "test",
//			wantErr:   errors.New("err do request"),
//			hashKey:   "some",
//		},
//		{
//			name:      "Empty hash key",
//			token:     tok,
//			mock:      func(s *MockSendInterface) {},
//			testEvent: "test",
//			wantErr:   errors.New("Empty hash key "),
//			address:   "localhost:8000",
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			s := NewMockSendInterface(c)
//			test.mock(s)
//
//			sss := &Sender{token: test.token, hashKey: test.hashKey, client: http.DefaultClient,
//				address: test.address}
//
//			if _, err := sss.GetDataRequest(test.testEvent); (err != nil) != (test.wantErr != nil) {
//				t.Errorf("GetDataRequest() error = %v, wantErr %v", err, test.wantErr != nil)
//			}
//		})
//	}
//}
//
//func TestSenderDeleteData(t *testing.T) {
//	//const tok = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE3Njc2NTMsIlVzZXJJRCI6M30.L_nb24gCL163FFNJsfgvIsqY3u_dho2f6VdJLySEVe8"
//	type mockBehaviour func(s *MockSendInterface)
//	tests := []struct {
//		name      string
//		token     string
//		mock      mockBehaviour
//		testEvent string
//		wantErr   error
//		hashKey   string
//		address   string
//	}{
//		{
//			name:      "Ok",
//			token:     tok,
//			mock:      func(s *MockSendInterface) {},
//			testEvent: "testing",
//			wantErr:   nil,
//			hashKey:   "some",
//			address:   "localhost:8000",
//		},
//		{
//			name:      "Empty token",
//			mock:      func(s *MockSendInterface) {},
//			testEvent: "testing",
//			wantErr:   errors.New("Token is empty, try to login "),
//			hashKey:   "some",
//			address:   "localhost:8000",
//		},
//		{
//			name:      "err to do request",
//			token:     tok,
//			mock:      func(s *MockSendInterface) {},
//			testEvent: "testing",
//			wantErr:   errors.New("err to do request"),
//			hashKey:   "some",
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			s := NewMockSendInterface(c)
//			test.mock(s)
//
//			sss := &Sender{token: test.token, hashKey: test.hashKey, client: http.DefaultClient,
//				address: test.address}
//
//			if err := sss.DeleteDataRequest(test.testEvent); (err != nil) != (test.wantErr != nil) {
//				t.Errorf("DeleteDataRequest() error = %v, wantErr %v", err, test.wantErr != nil)
//			}
//		})
//	}
//}
