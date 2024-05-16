package handler

import (
	"io"
	"net/http"

	_ "password_keeper/cmd/server/docs"

	"github.com/go-chi/chi/v5"
)

// @Summary SetData
// @Tags data
// @Description set data
// @ID set-data
// @Accept text
// @Produce text
// @Param input text
// @Success 201 {string} added
// @Failure 400 {string} err to get body or err to unmarshal body
// @Failure 500 {string} err to set data
// @Router api/set/{event} [post]

func (h *Handler) SetData(rw http.ResponseWriter, r *http.Request) {
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

// @Summary GetData
// @Tags data
// @Description get data by event
// @ID get-data
// @Produce text
// @Param {event} in url
// @Success 200 {string} data
// @Failure 500 {string} err to get data
// @Router api/get/{event} [get]
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

// @Summary DeleteData
// @Tags data
// @Description delete data by event
// @ID delete-data
// @Param {event} in url
// @Success 204 success
// @Failure 500 {string} err to delete data
// @Router api/delete/{event} [get]

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
