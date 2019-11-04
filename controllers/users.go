package controllers

import (
	"database/sql"
	"encoding/json"
	"essentials/models"
	"essentials/payloads/response"
	"log"
	"net/http"
)

// Users : struct for set Users Dependency Injection
type Users struct {
	Db *sql.DB
}

// List : http handler for returning list of users
func (u *Users) List(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	list, err := user.List(u.Db)
	if err != nil {
		log.Println("get user list", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var respList []response.UserResponse
	for _, l := range list {
		var resp response.UserResponse
		resp.Transform(&l)
		respList = append(respList, resp)
	}

	data, err := json.Marshal(respList)
	if err != nil {
		log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}
}
