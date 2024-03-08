package controllers

import (
	"encoding/json"
	m "modul5/models"
	"net/http"
)

// Send Error Umum
func sendErrorResponse(w http.ResponseWriter, status int, message string) {
	var response m.ErrorResponse
	response.Status = status
	response.Message = message

	// Mengirimkan response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Send Success Umum
func sendSuccessResponse(w http.ResponseWriter, status int, message string) {
	var response m.SuccessResponse
	response.Status = status
	response.Message = message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Send Succes Khusus Insert User
func sendGetUsersResponse(w http.ResponseWriter, status int, message string, data []m.Users) {
	var response m.UsersResponse
	response.Data = data
	response.Status = status
	response.Message = message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Send Succes Khusus Insert Product
func sendGetProductResponse(w http.ResponseWriter, status int, message string) {
	var response m.ProductsResponse
	response.Status = status
	response.Message = message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
