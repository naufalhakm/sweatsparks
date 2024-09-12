package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sweatsparks/internal/commons/response"
	"sweatsparks/internal/params"
	"sweatsparks/internal/services"

	"github.com/gorilla/mux"
)

type ProfileController interface {
	CreateProfile(w http.ResponseWriter, r *http.Request)
	GetDetailProfile(w http.ResponseWriter, r *http.Request)
	GetAllProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
}

type ProfileControllerImpl struct {
	ProfileService services.ProfileService
}

func NewProfileController(profileService services.ProfileService) ProfileController {
	return &ProfileControllerImpl{
		ProfileService: profileService,
	}
}

func (controller *ProfileControllerImpl) CreateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req params.ProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := response.BadRequestError("Invalid input")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err := controller.ProfileService.CreateProfileUser(r.Context(), &req)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.CreatedSuccessWithPayload("Success create profile user")
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *ProfileControllerImpl) GetDetailProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, _ := strconv.Atoi(userIDStr)
	result, err := controller.ProfileService.GetProfileUser(r.Context(), userID)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get data detail profile", result)
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *ProfileControllerImpl) GetAllProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, _ := strconv.Atoi(userIDStr)

	location := vars["location"]

	gender := vars["gender"]
	result, err := controller.ProfileService.GetAllProfileUser(r.Context(), userID, gender, location)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get all detail profile", result)
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *ProfileControllerImpl) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req params.ProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := response.BadRequestError("Invalid input")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err := controller.ProfileService.UpdateProfileUser(r.Context(), &req)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.CreatedSuccessWithPayload("Success update profile user")
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}
