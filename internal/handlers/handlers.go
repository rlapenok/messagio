package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rlapenok/messagio/internal/models"
	"github.com/rlapenok/messagio/internal/state"
)

func SendMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data models.Message
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := state.State.SaveMessage(ctx, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Your message has been applied and will be processed"))

}
func GetSats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	ctx := r.Context()
	data, err := state.State.GetSats(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
