package service

import (
	"context"
	"errors"
	"testing"

	_ "github.com/lib/pq"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"password_keeper/config/server"
	"password_keeper/internal/server/repository"
)

func TestSetData(t *testing.T) {
	type mockBehaviour func(s *MockDataServiceInterface)
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		id            int
		testData      []byte
		testEvent     string
		wantErr       error
	}{
		{
			name:          "success",
			mockBehaviour: func(s *MockDataServiceInterface) {},
			wantErr:       nil,
			testEvent:     "testEvent",
			testData:      []byte("testData"),
			id:            1,
		},
		{
			name:          "err to do exec",
			mockBehaviour: func(s *MockDataServiceInterface) {},
			wantErr:       errors.New("test error"),
			testEvent:     "testEvent",
			testData:      []byte("testData"),
			id:            9999999999,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := server.NewServConfig(
				server.WithDBAddress("postgres"),
				server.WithDBPort("5444"),
				server.WithDBUsername("postgres"),
				server.WithDBPassword("qwerty"),
				server.WithDBDatabase("postgres"),
			)

			db, err := repository.InitDataBase(ctx, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBUsername, cfg.DBPassword)
			require.NoError(t, err)

			defer db.Close()

			rep := repository.NewRepository(db)

			data := NewMockDataServiceInterface(c)
			test.mockBehaviour(data)

			service := NewDataService(rep)

			if err = service.SetData(test.id, test.testData, test.testEvent); (err != nil) != (test.wantErr != nil) {
				t.Errorf("SetData error = %v, wantErr %v", err, test.wantErr != nil)
			}

		})
	}
}

func TestGetData(t *testing.T) {
	type mockBehaviour func(s *MockDataServiceInterface)
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		id            int
		testEvent     string
		wantErr       error
	}{
		{
			name:          "success",
			mockBehaviour: func(s *MockDataServiceInterface) {},
			wantErr:       nil,
			testEvent:     "testEvent",
			id:            1,
		},
		{
			name:          "success",
			mockBehaviour: func(s *MockDataServiceInterface) {},
			wantErr:       errors.New("test error"),
			id:            1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := server.NewServConfig(
				server.WithDBAddress("postgres"),
				server.WithDBPort("5444"),
				server.WithDBUsername("postgres"),
				server.WithDBPassword("qwerty"),
				server.WithDBDatabase("postgres"),
			)

			db, err := repository.InitDataBase(ctx, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBUsername, cfg.DBPassword)
			require.NoError(t, err)

			defer db.Close()

			rep := repository.NewRepository(db)

			data := NewMockDataServiceInterface(c)
			test.mockBehaviour(data)

			service := NewDataService(rep)

			if _, err = service.GetData(test.id, test.testEvent); (err != nil) != (test.wantErr != nil) {
				t.Errorf("GetData() error = %v, wantErr %v", err, test.wantErr != nil)
			}

		})
	}
}

func TestDeleteData(t *testing.T) {
	type mockBehaviour func(s *MockDataServiceInterface)
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		id            int
		testEvent     string
		wantErr       error
	}{
		{
			name:          "success",
			mockBehaviour: func(s *MockDataServiceInterface) {},
			wantErr:       nil,
			testEvent:     "testID",
			id:            15,
		},
		{
			name:          "err to delete",
			mockBehaviour: func(s *MockDataServiceInterface) {},
			wantErr:       nil,
			testEvent:     "testID",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := server.NewServConfig(
				server.WithDBAddress("postgres"),
				server.WithDBPort("5444"),
				server.WithDBUsername("postgres"),
				server.WithDBPassword("qwerty"),
				server.WithDBDatabase("postgres"),
			)

			db, err := repository.InitDataBase(ctx, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBUsername, cfg.DBPassword)
			require.NoError(t, err)

			defer db.Close()

			rep := repository.NewRepository(db)

			data := NewMockDataServiceInterface(c)
			test.mockBehaviour(data)

			service := NewDataService(rep)

			if err = service.DeleteData(test.id, test.testEvent); (err != nil) != (test.wantErr != nil) {
				t.Errorf("GetData() error = %v, wantErr %v", err, test.wantErr != nil)
			}

		})
	}
}
