package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handlers) PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := h.repo.MonitorDatabase()
		if err != nil {
			if h.logger != nil {
				h.logger.Errorf("PingHandler: database monitor error: %v", err)
			}
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		send(w, http.StatusOK, state)
	}
}

func bindJson(r io.Reader, obj any) error {
	return json.NewDecoder(r).Decode(obj)
}

func send(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set(ContentType, AppJson)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

const (
	ContentType = "Content-Type"
	AppJson     = "application/json"
)
