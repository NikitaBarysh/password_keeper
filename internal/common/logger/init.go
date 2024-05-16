// Package logger - пакет в котором лежат методы и middleware логера
package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// Logger - структура, в которой лежит сущность Logger
type Logger struct {
	Log *zap.Logger
}

// InitLogger - создаем логгер со своими настройками
func InitLogger() *Logger {
	_, err := os.OpenFile("project_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	cfg := zap.NewProductionConfig()
	//cfg.OutputPaths = []string{"project_logs.log"}
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &Logger{Log: log}
}

// Info - метод, который выводит ошибки уровня Info
func (l *Logger) Info(msg string, opt ...any) {
	l.Log.Sugar().Infof(fmt.Sprintf(msg, opt...))
}

// Error - метод, который выводит ошибки уровня Error
func (l *Logger) Error(msg string, opt ...any) {
	l.Log.Sugar().Errorf(fmt.Sprintf(msg, opt...))
}

// Fatal - метод, который выводит ошибки уровня Fatal
func (l *Logger) Fatal(msg string, opt ...any) {
	l.Log.Sugar().Fatal(fmt.Sprintf(msg, opt...))
}
