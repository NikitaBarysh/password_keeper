// Package service - пакет в котором реализовано безнес-логика сервера
package service

import (
	"fmt"

	"password_keeper/internal/server/repository"
)

// DataService - структура в которой есть зависимость от repository
type DataService struct {
	rep *repository.Repository
}

// NewDataService - создаем структуру
func NewDataService(newRep *repository.Repository) *DataService {
	return &DataService{
		rep: newRep,
	}
}

// SetData - метод добавления новых данных в репозиторий
func (s *DataService) SetData(id int, text []byte, eventType string) error {
	if err := s.rep.SetRepData(id, text, eventType); err != nil {
		return fmt.Errorf("SetData: %w", err)
	}
	return nil
}

// GetData - метод получения данных из репозитория
func (s *DataService) GetData(id int, eventType string) ([]byte, error) {
	data, err := s.rep.GetRepData(id, eventType)
	if err != nil {
		return nil, fmt.Errorf("GetData: %w", err)
	}
	return data, nil
}

// DeleteData - метод удаления данных из репозитория
func (s *DataService) DeleteData(id int, eventType string) error {
	err := s.rep.DeleteRepData(id, eventType)
	if err != nil {
		return fmt.Errorf("DeleteData: %w", err)
	}
	return nil
}
