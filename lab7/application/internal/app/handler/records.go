package handler

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

func (h *Handlers) GetAllRecordsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var records []string

		var cursor uint64
		count := int64(64)

		for {
			keys, nextCursor, err := h.cache.Scan(r.Context(), cursor, "*", count)
			if err != nil {
				h.logger.Errorf("Failed to scan cache: %v", err)
				return
			}

			for _, key := range keys {
				val, err := h.cache.Get(r.Context(), key)
				if err == redis.Nil {
					continue
				}
				if err != nil {
					h.logger.Errorf("Failed to get value from cache: %v", err)
					continue
				}

				valStr, ok := val.(string)
				if !ok {
					continue
				}

				records = append(records, valStr)
			}

			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}

		if records == nil {
			records = []string{}
		}

		data := map[string]any{
			"records": records,
			"count": len(records),
		}
		
		w.Header().Set("Content-Type", "application/json")
		send(w, http.StatusOK, data)
	}
}