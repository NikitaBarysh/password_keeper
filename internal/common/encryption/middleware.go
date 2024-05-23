// Package encryption - пакет, в котором происходит шифрование данных
package encryption

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// DecryptMiddleware - middleware, который расшифровывает входящие данные авторизации по ручке
func DecryptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if dataDecryptor != nil {
			buf, _ := io.ReadAll(r.Body)

			msg, err := dataDecryptor.Decrypt(buf)
			if err != nil {
				log.Printf("DecryptMiddleware: error decrypting message: %s", err)
				http.Error(rw, "DecryptMiddleware: error decrypting message", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(msg))

		}
		h.ServeHTTP(rw, r)
	})
}
