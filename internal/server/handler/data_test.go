package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/server/service"
)

func TestSetData(t *testing.T) {
	type mockBehavior func(s *service.MockDataServiceInterface, inputData entity.Data)
	tests := []struct {
		name         string
		inputBody    string
		inputData    entity.Data
		mockBehavior mockBehavior
		expectedCode int
		param        map[string]any
	}{
		{
			name:      "Test #1 successful add",
			inputBody: `test payload`,
			inputData: entity.Data{
				Payload:   []byte("test payload"),
				EventType: "testEvent",
			},
			mockBehavior: func(s *service.MockDataServiceInterface, inputData entity.Data) {
				s.EXPECT().SetData(1, inputData.Payload, inputData.EventType).Return(nil)
			},
			expectedCode: http.StatusCreated,
			param:        map[string]any{"event": "testEvent"},
		},
		{
			name:      "Test #2 err to add",
			inputBody: `test payload`,
			inputData: entity.Data{
				Payload:   []byte("test payload"),
				EventType: "testEvent",
			},
			mockBehavior: func(s *service.MockDataServiceInterface, inputData entity.Data) {
				s.EXPECT().SetData(1, inputData.Payload, inputData.EventType).Return(errors.New("test error"))
			},
			expectedCode: http.StatusInternalServerError,
			param:        map[string]any{"event": "testEvent"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest(http.MethodPost, "/api/set/test", bytes.NewBufferString(test.inputBody))
			rctx := chi.NewRouteContext()
			for k, v := range test.param {
				strval := v.(string)
				rctx.URLParams.Add(k, strval)
			}
			rw := httptest.NewRecorder()

			ctx := context.WithValue(context.Background(), "user", 1)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

			data := service.NewMockDataServiceInterface(c)
			test.mockBehavior(data, test.inputData)
			serv := &service.Service{DataServiceInterface: data}

			handler := NewHandler(serv)
			handler.SetData(rw, req.WithContext(ctx))

			res := rw.Result()

			assert.Equal(t, test.expectedCode, res.StatusCode, test.name)
		})
	}
}

func TestGetData(t *testing.T) {
	type mockBehaviour func(s *service.MockDataServiceInterface, eventType string)

	tests := []struct {
		name          string
		eventType     string
		mockBehaviour mockBehaviour
		expectedCode  int
		param         map[string]any
	}{
		{
			name:      "Test #1 get data",
			eventType: "testEvent",
			mockBehaviour: func(s *service.MockDataServiceInterface, eventType string) {
				s.EXPECT().GetData(1, eventType).Return([]byte("test payload"), nil)
			},
			expectedCode: http.StatusOK,
			param:        map[string]any{"event": "testEvent"},
		},
		{
			name:      "Test #2 err to get",
			eventType: "testEvent",
			mockBehaviour: func(s *service.MockDataServiceInterface, eventType string) {
				s.EXPECT().GetData(1, eventType).Return(nil, errors.New("test error"))
			},
			expectedCode: http.StatusInternalServerError,
			param:        map[string]any{"event": "testEvent"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest(http.MethodGet, "/api/get/testEvent", nil)
			rctx := chi.NewRouteContext()
			for k, v := range test.param {
				strval := v.(string)
				rctx.URLParams.Add(k, strval)
			}
			rw := httptest.NewRecorder()

			ctx := context.WithValue(context.Background(), "user", 1)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

			data := service.NewMockDataServiceInterface(c)
			test.mockBehaviour(data, test.eventType)
			serv := &service.Service{DataServiceInterface: data}

			handler := NewHandler(serv)
			handler.getData(rw, req.WithContext(ctx))

			res := rw.Result()

			assert.Equal(t, test.expectedCode, res.StatusCode, test.name)
		})

	}
}

func TestDeleteData(t *testing.T) {
	type mockBehaviour func(s *service.MockDataServiceInterface, eventType string)
	tests := []struct {
		name          string
		eventType     string
		mockBehaviour mockBehaviour
		expectedCode  int
		param         map[string]any
	}{
		{
			name:      "Test #1 delete data",
			eventType: "testEvent",
			mockBehaviour: func(s *service.MockDataServiceInterface, eventType string) {
				s.EXPECT().DeleteData(1, eventType).Return(nil)
			},
			expectedCode: http.StatusNoContent,
			param:        map[string]any{"event": "testEvent"},
		},
		{
			name:      "Test #2 err to delete",
			eventType: "testEvent",
			mockBehaviour: func(s *service.MockDataServiceInterface, eventType string) {
				s.EXPECT().DeleteData(1, eventType).Return(errors.New("test error"))
			},
			expectedCode: http.StatusInternalServerError,
			param:        map[string]any{"event": "testEvent"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest(http.MethodDelete, "/api/delete/testEvent", nil)
			rctx := chi.NewRouteContext()
			for k, v := range test.param {
				strval := v.(string)
				rctx.URLParams.Add(k, strval)
			}
			rw := httptest.NewRecorder()

			ctx := context.WithValue(context.Background(), "user", 1)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

			data := service.NewMockDataServiceInterface(c)
			test.mockBehaviour(data, test.eventType)
			serv := &service.Service{DataServiceInterface: data}
			handler := NewHandler(serv)
			handler.deleteData(rw, req.WithContext(ctx))
			res := rw.Result()
			assert.Equal(t, test.expectedCode, res.StatusCode, test.name)
		})
	}
}
