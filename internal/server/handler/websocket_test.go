package handler

//func TestWebsocket(t *testing.T) {
//	type mockBehaviour func(s *service.MockDataServiceInterface, inputData entity.WebDataMSG)
//	tests := []struct {
//		name          string
//		mockBehaviour mockBehaviour
//		input         entity.WebDataMSG
//		statusCode    int
//	}{
//		{
//			name: "test #1",
//			mockBehaviour: func(s *service.MockDataServiceInterface, inputData entity.WebDataMSG) {
//
//			},
//			input: entity.WebDataMSG{
//				Action:    "add",
//				EventType: "testEvent",
//				Payload:   []byte(`{"test":"test"}`),
//			},
//			statusCode: http.StatusOK,
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			data := service.NewMockDataServiceInterface(c)
//			test.mockBehaviour(data, test.input)
//
//			serv := &service.Service{DataServiceInterface: data}
//			handler := NewHandler(serv)
//
//			//req := httptest.NewRequest(http.MethodPost, "/ws", bytes.NewBuffer(test.input.Payload))
//			//w := httptest.NewRecorder()
//			//ctx := context.WithValue(context.Background(), "user", 1)
//
//			// Не знаю как передать контекст в запрос
//			s := httptest.NewServer(http.HandlerFunc(handler.handleSetWebsocket))
//
//			u := "ws" + strings.TrimPrefix(s.URL, "http")
//
//			ws, _, err := websocket.DefaultDialer.Dial(u, nil)
//			require.NoError(t, err)
//			defer ws.Close()
//
//			b, err := json.Marshal(test.input)
//			require.NoError(t, err)
//
//			ws.WriteMessage(websocket.TextMessage, b)
//		})
//	}
//}
