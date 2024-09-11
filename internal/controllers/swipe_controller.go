package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sweatsparks/internal/commons/response"
	"sweatsparks/internal/params"
	"sweatsparks/internal/services"

	"github.com/gorilla/mux"
)

type SwipeController interface {
	CreateSwipe(w http.ResponseWriter, r *http.Request)
	GetSwipeDetail(w http.ResponseWriter, r *http.Request)
	GetSwipeAll(w http.ResponseWriter, r *http.Request)
}

type SwipeControllerImpl struct {
	SwipeService services.SwipeService
}

func NewSwipeController(swipeService services.SwipeService) SwipeController {
	return &SwipeControllerImpl{
		SwipeService: swipeService,
	}
}

func (controller *SwipeControllerImpl) CreateSwipe(w http.ResponseWriter, r *http.Request) {
	var req params.SwipeRequest
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := response.BadRequestError("Invalid input")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err := controller.SwipeService.CreateSwipe(context.Background(), &req)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.CreatedSuccessWithPayload("Success insert swipe")

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *SwipeControllerImpl) GetSwipeDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	swiperIDStr := vars["swiperID"]
	swiperID, _ := strconv.Atoi(swiperIDStr)

	swipeeIDStr := vars["swiperID"]
	swipeeID, _ := strconv.Atoi(swipeeIDStr)

	result, err := controller.SwipeService.GetSwipeBySwiperAndSwipee(context.Background(), swiperID, swipeeID)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get detail data swipe", result)

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *SwipeControllerImpl) GetSwipeAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	swipeeIDStr := vars["swiperID"]
	swipeeID, _ := strconv.Atoi(swipeeIDStr)

	result, err := controller.SwipeService.GetAllSwipeeNotMatchBySwipee(context.Background(), swipeeID)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get all data swipe", result)

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}
