package encryption

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func DecryptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if dataDecryptor != nil {
			buf, _ := io.ReadAll(r.Body)

			msg, err := dataDecryptor.Decrypt(buf)
			if err != nil {
				log.Printf("Error decrypting message: %s", err)
				http.Error(rw, "Error decrypting message", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(msg))

		}
		h.ServeHTTP(rw, r)
	})
}
