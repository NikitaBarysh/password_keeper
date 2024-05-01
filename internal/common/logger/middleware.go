package logger

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		size   int
		status int
		head   string
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size += size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

func (l *loggingResponseWriter) Header() http.Header {
	return l.ResponseWriter.Header()
}

func (l *loggingResponseWriter) Hijacking() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := l.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func (l *loggingResponseWriter) Flush() {
	if f, ok := l.ResponseWriter.(http.Flusher); ok {
		if l.ResponseWriter. == 0 {
			l.responseData.status = 200
		}
		f.Flush()
	}
}

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		response := &responseData{
			size:   0,
			status: 0,
			head:   "",
		}
		lw := loggingResponseWriter{
			ResponseWriter: rw,
			responseData:   response,
		}

		handler.ServeHTTP(&lw, r)
		duration := time.Since(start)

		InitLogger().Info("",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", response.status),
			zap.Duration("duration", duration),
			zap.Int("size", response.size),
			zap.String("header", response.head))
	})
}
