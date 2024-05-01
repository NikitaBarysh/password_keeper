package service

import (
	"context"
	"fmt"

	"password_keeper/internal/common/logger"
	"password_keeper/internal/server/repository"
)

type DataService struct {
	rep     *repository.Repository
	logging *logger.Logger
}

func NewDataService(newRep *repository.Repository, logger *logger.Logger) *DataService {
	return &DataService{
		rep:     newRep,
		logging: logger,
	}
}

func (s *DataService) SetData(ctx context.Context, id int, text []byte, eventType string) error {
	if err := s.rep.SetRepData(ctx, id, text, eventType); err != nil {
		s.logging.Error("err to set data in DB: ", err)
		return fmt.Errorf("err to set data in DB: %w", err)
	}
	return nil
}

func (s *DataService) GetData(ctx context.Context, id int, eventType string) ([]byte, error) {
	data, err := s.rep.GetRepData(ctx, id, eventType)
	if err != nil {
		s.logging.Error("err to get data from DB: ", err)
		return nil, fmt.Errorf("err to get data from DB: %w", err)
	}
	return data, nil
}

func (s *DataService) DeleteData(ctx context.Context, id int, eventType string) error {
	err := s.rep.DeleteRepData(ctx, id, eventType)
	if err != nil {
		s.logging.Error("err to delete data from DB: ", err)
		return fmt.Errorf("err to delete data from DB: %w", err)
	}
	return nil
}
