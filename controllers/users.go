package controllers

import (
	"database/sql"
	"essentials/libraries/api"
	"essentials/models"
	"essentials/payloads/request"
	"essentials/payloads/response"
	"essentials/usecases"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Users : struct for set Users Dependency Injection
type Users struct {
	Db  *sql.DB
	Log *log.Logger
}

// List : http handler for returning list of users
func (u *Users) List(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	list, err := user.List(u.Db)
	if err != nil {
		u.Log.Println("get user list", err)
		api.ResponseError(w, err)
		return
	}

	var respList []response.UserResponse
	for _, l := range list {
		var resp response.UserResponse
		resp.Transform(&l)
		respList = append(respList, resp)
	}

	api.ResponseOK(w, respList, http.StatusOK)
}

// Create new user
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	uc := usecases.UserUsecase{Log: u.Log, Db: u.Db}
	data, err := uc.Create(r)
	if err != nil {
		api.ResponseError(w, err)
		return
	}

	api.ResponseOK(w, data, http.StatusOK)
}

// View user by id
func (u *Users) View(w http.ResponseWriter, r *http.Request) {
	paramID := r.Context().Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Println("convert param to id", err)
		api.ResponseError(w, err)
		return
	}

	user := new(models.User)
	user.ID = uint64(id)
	err = user.Get(u.Db)
	if err != nil {
		u.Log.Println("Get User", err)
		api.ResponseError(w, err)
		return
	}

	resp := new(response.UserResponse)
	resp.Transform(user)
	api.ResponseOK(w, resp, http.StatusOK)
}

// Update user by id
func (u *Users) Update(w http.ResponseWriter, r *http.Request) {
	paramID := r.Context().Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Println("convert param to id", err)
		api.ResponseError(w, err)
		return
	}

	user := new(models.User)
	user.ID = uint64(id)
	err = user.Get(u.Db)
	if err != nil {
		u.Log.Println("Get User", err)
		api.ResponseError(w, err)
		return
	}

	userRequest := new(request.UserRequest)
	err = api.Decode(r, &userRequest)
	if err != nil {
		u.Log.Printf("error decode user: %s", err)
		api.ResponseError(w, err)
		return
	}

	userUpdate := userRequest.Transform(user)
	err = userUpdate.Update(u.Db)
	if err != nil {
		u.Log.Printf("error update user: %s", err)
		api.ResponseError(w, err)
		return
	}

	resp := new(response.UserResponse)
	resp.Transform(userUpdate)
	api.ResponseOK(w, resp, http.StatusOK)
}

// Delete user by id
func (u *Users) Delete(w http.ResponseWriter, r *http.Request) {
	paramID := r.Context().Value(api.Ctx("ps")).(httprouter.Params).ByName("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.Log.Println("convert param to id", err)
		api.ResponseError(w, err)
		return
	}

	user := new(models.User)
	user.ID = uint64(id)
	err = user.Get(u.Db)
	if err != nil {
		u.Log.Println("Get User", err)
		api.ResponseError(w, err)
		return
	}

	err = user.Delete(u.Db)
	if err != nil {
		u.Log.Println("Delete User", err)
		api.ResponseError(w, err)
		return
	}

	api.ResponseOK(w, nil, http.StatusOK)
}
