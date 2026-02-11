package handler

import "net/http"

func (h *Handlers) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("good"))
		w.WriteHeader(http.StatusOK)
	}
}