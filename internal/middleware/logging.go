package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// LogRequestBody logs the request body and restores it for downstream handlers
func LogRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only log for methods that might have a body
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			if r.Body == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Read the body
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			// Log the body if present
			if len(bodyBytes) > 0 {
				log.Printf("Request Body: %s", string(bodyBytes))
			}

			// Restore the body
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		next.ServeHTTP(w, r)
	})
}
