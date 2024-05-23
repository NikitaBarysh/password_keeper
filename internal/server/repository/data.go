// Package repository - пакет в котором сохраняем данные в базу данных
package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DataRepository - структура, в которой лежит сущность ба
type DataRepository struct {
	db *sqlx.DB
}

// NewDataRepository - создает DataRepository
func NewDataRepository(newDB *sqlx.DB) *DataRepository {
	return &DataRepository{
		db: newDB,
	}
}

// SetRepData - вставляем данные пользователя
func (r *DataRepository) SetRepData(id int, text []byte, eventType string) error {
	_, err := r.db.Exec(`INSERT INTO data (user_id, data, event_type) VALUES ($1, $2, $3)`, id, text, eventType)
	if err != nil {
		return fmt.Errorf("SetRepData: err to do exec in DB: %w", err)
	}

	return nil
}

// GetRepData - получаем данные пользователя
func (r *DataRepository) GetRepData(id int, eventType string) ([]byte, error) {
	var data []byte

	row := r.db.QueryRow(`SELECT data FROM data where user_id = $1 AND event_Type = $2`, id, eventType)

	err := row.Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("GetRepData: err to do scan: %w", err)
	}

	return data, nil
}

// DeleteRepData - удаляем данные пользователя
func (r *DataRepository) DeleteRepData(id int, eventType string) error {
	_, err := r.db.Exec(`DELETE FROM data WHERE user_id = $1 and event_type = $2`, id, eventType)
	if err != nil {
		return fmt.Errorf("GetRepData: err to do exec in DB: %w", err)
	}
	return nil
}
