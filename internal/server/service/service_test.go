package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/server/repository"
)

func TestServiceGetData(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *repository.MockAuthorizationRepository, user entity.User)
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		user          entity.User
	}{
		{
			name: "OK",
			mockBehaviour: func(ctx context.Context, s *repository.MockAuthorizationRepository, user entity.User) {
				s.EXPECT().SetUserDB(ctx, user).Return(1, nil)
			},
			user: entity.User{
				Login:    "admin",
				Password: "admin",
			},
		},
		{
			name: "No ok",
			mockBehaviour: func(ctx context.Context, s *repository.MockAuthorizationRepository, user entity.User) {
				s.EXPECT().SetUserDB(ctx, user).Return(0, errors.New("err"))
			},
			user: entity.User{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := repository.NewMockAuthorizationRepository(c)
			test.mockBehaviour(ctx, auth, test.user)
			r := &repository.Repository{AuthorizationRepository: auth}

			r.SetUserDB(ctx, test.user)
		})
	}
}
