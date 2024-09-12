package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sweatsparks/internal/commons/response"
	"sweatsparks/internal/services"

	"github.com/gorilla/mux"
)

type MessageController interface {
	GetMessageByMatchID(w http.ResponseWriter, r *http.Request)
}

type MessageControllerImpl struct {
	MessageService services.MessageService
}

func NewMessageController(messageService services.MessageService) MessageController {
	return &MessageControllerImpl{
		MessageService: messageService,
	}
}

func (controller *MessageControllerImpl) GetMessageByMatchID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	matchIDStr := vars["matchID"]
	matchID, _ := strconv.Atoi(matchIDStr)

	message, err := controller.MessageService.GetMessageByMatchId(r.Context(), matchID)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get data message", message)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
