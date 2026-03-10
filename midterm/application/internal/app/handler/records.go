package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) GetAllRecordsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := h.repo.GetAllRecords(r.Context())
		if err != nil {
			if h.logger != nil {
				h.logger.Errorf("GetAllRecordsHandler: error fetching records: %v", err)
			}
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		send(w, http.StatusOK, records)
	}
}

func (h *Handlers) GetRecordByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/records/")
		
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		record, err := h.repo.GetRecordByID(r.Context(), idInt)
		if err != nil {
			if h.logger != nil {
				h.logger.Errorf("GetRecordByIDHandler: error fetching record by ID: %v", err)
			}
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if record == nil {
			w.Write([]byte("record not found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		send(w, http.StatusOK, record)
	}
}

func (h *Handlers) InsertRecordHandler() http.HandlerFunc {
	type request struct {
		Payload string `json:"payload"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := bindJson(r.Body, &req); err != nil {
			if h.logger != nil {
				h.logger.Errorf("InsertRecordHandler: error parsing request body: %v", err)
			}
			w.Write([]byte("invalid request body"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := h.repo.InsertRecord(r.Context(), req.Payload)
		if err != nil {
			if h.logger != nil {
				h.logger.Errorf("InsertRecordHandler: error inserting record: %v", err)
			}
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		send(w, http.StatusCreated, map[string]int{"id": id})
	}
}