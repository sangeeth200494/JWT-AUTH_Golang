package userhandlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sangeeth200494/JWT-AUTH_Golang/database"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
)

func RetrieveUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 400, Message: "invalid user id", Details: err.Error()})
		return
	}

	db, err := database.DBConnection()
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "database connection failed", Details: err.Error()})
		return
	}
	defer database.DBC(db)

	var user models.User
	result := db.Select("id, username, email, full_name, phone_number, status, role, created_at, last_login").First(&user, id)
	if result.Error != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 404, Message: "user not found", Details: result.Error.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "user retrieved successfully",
		"data":    user,
	})
}
