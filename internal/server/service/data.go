package service

import (
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

func (s *DataService) SetData(id int, text []byte, eventType string) error {
	if err := s.rep.SetRepData(id, text, eventType); err != nil {
		return fmt.Errorf("err to set data in DB: %w", err)
	}
	return nil
}

func (s *DataService) GetData(id int, eventType string) ([]byte, error) {
	data, err := s.rep.GetRepData(id, eventType)
	if err != nil {
		return nil, fmt.Errorf("err to get data from DB: %w", err)
	}
	return data, nil
}

func (s *DataService) DeleteData(id int, eventType string) error {
	err := s.rep.DeleteRepData(id, eventType)
	if err != nil {
		return fmt.Errorf("err to delete data from DB: %w", err)
	}
	return nil
}
