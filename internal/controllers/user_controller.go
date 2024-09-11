package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"sweatsparks/internal/commons/response"
	"sweatsparks/internal/params"
	"sweatsparks/internal/services"
)

type UserController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
}

type UserControllerImpl struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var req params.UserRegisterRequest
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := response.BadRequestError("Invalid input")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err := controller.UserService.RegisterUser(context.Background(), &req)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.CreatedSuccessWithPayload("Success register new user")

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var req params.UserLoginRequest
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := response.BadRequestError("Invalid input")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, err := controller.UserService.LoginUser(context.Background(), &req)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success login user", user)

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func (controller *UserControllerImpl) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user, err := controller.UserService.GetAllUser(context.Background())
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.GeneralSuccessCustomMessageAndPayload("Success get all data users", user))
}
