package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"password_keeper/internal/common/logger"
)

type DataRepository struct {
	db      *sqlx.DB
	logging *logger.Logger
}

func NewDataRepository(newDB *sqlx.DB, logger *logger.Logger) *DataRepository {
	return &DataRepository{
		db:      newDB,
		logging: logger,
	}
}

func (r *DataRepository) SetRepData(id int, text []byte, eventType string) error {
	_, err := r.db.Exec(`INSERT INTO data (user_id, data, event_type) VALUES ($1, $2, $3)`, id, text, eventType)
	if err != nil {
		r.logging.Error("err to do exec in DB: ", err)
		return fmt.Errorf("err to do exec in DB: %w", err)
	}

	return nil
}

func (r *DataRepository) GetRepData(id int, eventType string) ([]byte, error) {
	var data []byte

	row := r.db.QueryRow(`SELECT data FROM data where user_id = $1 AND event_Type = $2`, id, eventType)

	err := row.Scan(&data)
	if err != nil {
		r.logging.Error("err to do scan: ", err)
		return nil, fmt.Errorf("err to do scan: %w", err)
	}

	return data, nil
}

func (r *DataRepository) DeleteRepData(id int, eventType string) error {
	_, err := r.db.Exec(`DELETE FROM data WHERE user_id = $1 and event_type = $2`, id, eventType)
	if err != nil {
		r.logging.Error("err to do exec in DB: ", err)
		return fmt.Errorf("err to do exec in DB: %w", err)
	}
	return nil
}
