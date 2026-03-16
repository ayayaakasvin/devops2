package middleware

import (
	"application-for-kubernetes/internal/domain"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const recordTTL = 15 * time.Minute

func (m *Middlewares) CacheWriteMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := writeRecord(m.cache, r); err != nil {
			m.logger.Errorf("Failed to write record to cache: %v", err)
		}

		next.ServeHTTP(w, r)
	}
}

func writeRecord(cache domain.Cache, r *http.Request) error {
	id, recordString := record(r)
	err := cache.Set(r.Context(), id, recordString, recordTTL)
	return err
}

// returns random record ID and record string format "<URL>:<IP_ADDRESS>:<TIME>"
func record(r *http.Request) (string, string) {
	id := uuid.NewString()
	recordString := fmt.Sprintf("%s:%s:%s", r.URL.String(), r.RemoteAddr, time.Now().Format(time.RFC3339))
	return id, recordString
}