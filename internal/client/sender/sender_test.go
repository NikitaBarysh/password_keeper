package sender

//import (
//	"net/http"
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/require"
//	"password_keeper/config/client"
//)
//
//func TestSender(t *testing.T) {
//	token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTQ3NDQwNTgsIlVzZXJJRCI6MTV9.SWIMnCXfBGgfzi8J0IjNSG9_gZDCcKtHEmXT5RR7wpA"
//	type mockBehaviour func(s *MockSendInterface, data string, event string)
//	tests := []struct {
//		name      string
//		mock      mockBehaviour
//		testData  string
//		testEvent string
//	}{
//		{
//			name: "Ok",
//			mock: func(s *MockSendInterface, data string, event string) {
//				s.EXPECT().PostDataRequest(data, event).Return(nil)
//			},
//			testData:  "hi",
//			testEvent: "test",
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			cfg := client.NewClient()
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			s := NewMockSendInterface(c)
//			test.mock(s, test.testData, test.testEvent)
//
//			sss := &Sender{sendInterface: s, token: token, cfg: cfg, client: http.DefaultClient}
//
//			err := sss.PostDataRequest("hi", "test")
//			require.NoError(t, err)
//		})
//	}
//}
