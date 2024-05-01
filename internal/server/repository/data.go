package repository

import (
	"context"
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

func (r *DataRepository) SetRepData(ctx context.Context, id int, text []byte, eventType string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO data (user_id, data, event_type) VALUES ($1, $2, $3)`, id, text, eventType)
	if err != nil {
		r.logging.Error("err to do exec in DB: ", err)
		return fmt.Errorf("err to do exec in DB: %w", err)
	}

	return nil
}

func (r *DataRepository) GetRepData(ctx context.Context, id int, eventType string) ([]byte, error) {
	var data []byte

	row := r.db.QueryRowContext(ctx, `SELECT data FROM data where user_id = $1 AND event_Type = $2`, id, eventType)

	err := row.Scan(&data)
	if err != nil {
		r.logging.Error("err to do scan: ", err)
		return nil, fmt.Errorf("err to do scan: %w", err)
	}

	//rows, err := r.db.QueryxContext(ctx, `SELECT * FROM data where user_id = $1`, id)
	//if err != nil {
	//	r.logging.Error("err to do query in DB: ", err)
	//	return nil, fmt.Errorf("err to do query in DB: %w", err)
	//}
	//
	//dataSlice := make([][]byte, 0)
	//
	//for rows.Next() {
	//	var data entity.Payload
	//	err = rows.Scan(&data.Id, &data.Payload)
	//	if err != nil {
	//		r.logging.Error("err to do row scan: ", err)
	//		return nil, fmt.Errorf("err to do row scan: %w", err)
	//	}
	//	dataSlice = append(dataSlice, data.Payload)
	//}
	//
	return data, nil
}

func (r *DataRepository) DeleteRepData(ctx context.Context, id int, eventType string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM data WHERE user_id = $1 and event_type = $2`, id, eventType)
	if err != nil {
		r.logging.Error("err to do exec in DB: ", err)
		return fmt.Errorf("err to do exec in DB: %w", err)
	}
	return nil
}
