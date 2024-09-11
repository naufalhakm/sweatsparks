package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sweatsparks/internal/commons/response"
	"sweatsparks/internal/middleware"
	"sweatsparks/internal/params"
	"sweatsparks/internal/services"

	"github.com/gorilla/mux"
)

type MatchController interface {
	CreateMatch(w http.ResponseWriter, r *http.Request)
	GetDetailMatchUser(w http.ResponseWriter, r *http.Request)
	GetAllMatchUser(w http.ResponseWriter, r *http.Request)
}

type MatchControllerImpl struct {
	MatchService services.MatchService
}

func NewMatchController(matchService services.MatchService) MatchController {
	return &MatchControllerImpl{
		MatchService: matchService,
	}
}

func (controller *MatchControllerImpl) CreateMatch(w http.ResponseWriter, r *http.Request) {
	var req params.MatchRequest
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := response.BadRequestError("Invalid input")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err := controller.MatchService.CreateMatchUser(context.Background(), &req)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.CreatedSuccessWithPayload("Success create new a match")

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *MatchControllerImpl) GetDetailMatchUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"errors": "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID2, _ := strconv.Atoi(userIDStr)

	match, err := controller.MatchService.FindMatchDetailByUserID(context.Background(), int(userID), userID2)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get detail match", match)

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *MatchControllerImpl) GetAllMatchUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"errors": "Unauthorized",
		})
		return
	}

	match, err := controller.MatchService.FindMatchAllByUserID(context.Background(), int(userID))
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get all match", match)

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}
