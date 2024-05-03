package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	Log *zap.Logger
}

func InitLogger() *Logger {
	_, err := os.OpenFile("project_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"project_log.log"}
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &Logger{Log: log}
}

func (l *Logger) Info(msg string, opt ...any) {
	l.Log.Sugar().Infow(fmt.Sprintf(msg, opt...))
}

func (l *Logger) Error(msg string, opt ...any) {
	l.Log.Sugar().Errorw(fmt.Sprintf(msg, opt...))
}

func (l *Logger) Fatal(msg string, opt ...any) {
	l.Log.Sugar().Fatalw(fmt.Sprintf(msg, opt...))
}
