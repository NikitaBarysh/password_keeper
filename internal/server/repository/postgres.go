// Package repository - пакет в котором сохраняем данные в базу данных
package repository

import (
	"context"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// InitDataBase - подключаемся к базе данных
func InitDataBase(ctx context.Context, addr, port, dbName, user, pass string) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", addr, port, user, dbName, pass)
	log.Println("url: ", url)
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		log.Println("err open: ", err)
		return nil, fmt.Errorf("InitDataBase: err to create DB: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Println("err ping: ", err)
		return nil, fmt.Errorf("InitDataBase: err to ping DB: %w", err)
	}

	log.Println("db work: ", db)
	return db, nil
}
