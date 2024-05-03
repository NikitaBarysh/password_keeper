package handler

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) setData(rw http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "event")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "err to get body", http.StatusBadRequest)
		return
	}

	currUser := r.Context().Value("user").(int)

	err = h.service.DataServiceInterface.SetData(currUser, b, eventType)
	if err != nil {
		http.Error(rw, "err to set data", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("added"))
}

func (h *Handler) getData(rw http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "event")
	currUser := r.Context().Value("user").(int)

	data, err := h.service.DataServiceInterface.GetData(currUser, eventType)
	if err != nil {
		http.Error(rw, "err to get data", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (h *Handler) deleteData(rw http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "event")
	currUser := r.Context().Value("user").(int)

	err := h.service.DataServiceInterface.DeleteData(currUser, eventType)
	if err != nil {
		http.Error(rw, "err to delete data", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
