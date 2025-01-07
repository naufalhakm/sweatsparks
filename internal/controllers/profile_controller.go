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
	GetAllRecomendationUsers(w http.ResponseWriter, r *http.Request)
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

func (controller *ProfileControllerImpl) GetAllRecomendationUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, _ := strconv.Atoi(userIDStr)

	pageStr := vars["page"]
	page, _ := strconv.Atoi(pageStr)

	limitStr := vars["limit"]
	limit, _ := strconv.Atoi(limitStr)

	pageNum := 1
	limitSize := 5

	if page > 0 {
		pageNum = page
	}

	if limit > 0 {
		limitSize = limit
	}

	result, pagination, err := controller.ProfileService.GetUserRecomendation(r.Context(), userID, pageNum, limitSize)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	type Response struct {
		Users      interface{} `json:"users"`
		Pagination interface{} `json:"pagination"`
	}

	var responses Response
	responses.Users = result
	responses.Pagination = pagination

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get all recomendation users", responses)
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
