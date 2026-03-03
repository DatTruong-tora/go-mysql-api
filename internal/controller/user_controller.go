package controller

import (
	"encoding/json"
	"go-mysql-api/internal/models"
	"go-mysql-api/internal/service"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type UserController struct {
	Service *service.UserService
}

func (c *UserController) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := c.Service.GetAllUsers()
		if err != nil {
			logrus.WithError(err).Error("Error in getting all users")
			c.sendError(w, http.StatusInternalServerError, "Không thể lấy danh sách")
			return
		}
		c.sendJSON(w, http.StatusOK, users, "Lấy danh sách thành công")

	case http.MethodPost:
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			logrus.WithError(err).Error("Error in decoding user")
			c.sendError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
			return
		}
		if err := c.Service.CreateUser(&u); err != nil {
			logrus.WithError(err).Error("Error in creating user")
			c.sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		c.sendJSON(w, http.StatusCreated, u, "Tạo thành công")
	}
}

func (c *UserController) HandleUserDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithError(err).Error("Error in converting user ID")
		c.sendError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, err := c.Service.GetUser(id)
		if err != nil {
			logrus.WithError(err).Error("Error in getting user")
			c.sendError(w, http.StatusNotFound, "User not found")
			return
		}
		c.sendJSON(w, http.StatusOK, u, "Success")

	case http.MethodPut:
		var u models.User
		json.NewDecoder(r.Body).Decode(&u)
		if err := c.Service.UpdateUser(id, &u); err != nil {
			logrus.WithError(err).Error("Error in updating user")
			c.sendError(w, http.StatusInternalServerError, "Update failed")
			return
		}
		c.sendJSON(w, http.StatusOK, nil, "Update success")

	case http.MethodDelete:
		if err := c.Service.RemoveUser(id); err != nil {
			logrus.WithError(err).Error("Error in deleting user")
			c.sendError(w, http.StatusInternalServerError, "Delete failed")
			return
		}
		c.sendJSON(w, http.StatusOK, nil, "Deleted success")
	}
}

func (c *UserController) sendJSON(w http.ResponseWriter, status int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.APIResponse{Success: true, Message: msg, Data: data})
}

func (c *UserController) sendError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.APIResponse{Success: false, Message: msg})
}
