package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	m "modul5/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// 1. GET ALL USERS
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	query := "SELECT * FROM users"
	// Read from Query Param
	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]
	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name='" + name[0] + "'  "
	}

	if age != nil {
		if name[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " age= '" + age[0] + "' "
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}
	var user m.Users
	var users []m.Users
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Password, &user.Email); err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "internal error")
			return
		} else {
			users = append(users, user)
		}
	}
	w.Header().Set("Content-Type", "application/json")

	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

// 1b. GET ALL USERS GORM
func GetAllUsers2(w http.ResponseWriter, r *http.Request) {
	db := gorm_connect(w)

	// Mengambil semua pengguna dari database
	var users []m.Users
	result := db.Find(&users) // bisa pake first
	if result.Error != nil {
		if result.RowsAffected == 0 {
			sendErrorResponse(w, 404, "Pengguna tidak ditemukan")
		} else {
			log.Println(result.Error)
			sendErrorResponse(w, 500, "Kesalahan internal server")
		}
		return
	}

	sendGetUsersResponse(w, 200, "Berhasil Get Data", users)
}

// 1c. GET USERS BY AGE
func GetUsersByAge(w http.ResponseWriter, r *http.Request) {
	db := gorm_connect(w)

	// Mengambil nilai parameter 'age' dari URL
	vars := mux.Vars(r)
	ageParam, ok := vars["age"]
	if !ok {
		sendErrorResponse(w, 500, "Parameter 'age' tidak ditemukan")
		return
	}

	// Contoh query raw untuk mengambil data berdasarkan kueri SQL
	query := "SELECT * FROM users WHERE age = ?"
	var users []m.Users
	result := db.Raw(query, ageParam).Scan(&users)

	if result.Error != nil {
		log.Println(result.Error)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Mengirimkan response JSON
	var response m.UsersResponse
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

// 2. INSERT USER (POST)
func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	query := "INSERT INTO users (ID, Name, Age, Address, Password, Email) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Mendapatkan data dari request body
	var user m.Users
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Menjalankan pernyataan Exec dengan data dari request body
	_, err = stmt.Exec(user.ID, user.Name, user.Age, user.Address)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User berhasil dibuat"))

	// Mengisi nilai objek users
	var users []m.Users
	users = append(users, user)

	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users

	// Mengirimkan response JSON
	json.NewEncoder(w).Encode(response)
}

// 2b. INSERT USER Pakai GORM
func InsertUser2(w http.ResponseWriter, r *http.Request) {
	// Connect gorm dan penangangn error db
	db := gorm_connect(w)

	// Mendapatkan data dari request body
	var user m.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 400, "Bad request, format data tidak sesuai")
		return
	}

	// Memasukkan data user ke database menggunakan GORM
	result := db.Create(&user)
	if result.Error != nil {
		log.Println(result.Error)
		sendErrorResponse(w, 500, "Gagal memasukan data user")
		return
	}

	// Respon json
	sendSuccessResponse(w, 201, "User berhasil dibuat")
}

// 3. UPDATE USER (PUT)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	// Mendapatkan ID dari path parameter
	vars := mux.Vars(r)
	userID := vars["id"]

	// Membaca data dari request body
	var updatedUser m.Users
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Mengeksekusi pernyataan SQL untuk mengupdate user berdasarkan ID
	query := "UPDATE users SET Name=?, Age=?, Address=? WHERE ID=?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	_, err = stmt.Exec(updatedUser.Name, updatedUser.Age, updatedUser.Address, userID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Mengeksekusi pernyataan SQL untuk mendapatkan data user yang telah diupdate
	querySelect := "SELECT ID, Name, Age, Address FROM users WHERE ID=?"
	rows, err := db.Query(querySelect, userID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}
	defer rows.Close()

	// Mengonversi hasil query ke dalam bentuk array users
	var users []m.Users
	for rows.Next() {
		var user m.Users
		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "internal error")
			return
		}
		users = append(users, user)
	}

	// Membuat response JSON
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Update success"
	response.Data = users

	// Mengirimkan response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 3b. UPDATE USER Pakai GORM
func UpdateUser2(w http.ResponseWriter, r *http.Request) {
	db := gorm_connect(w)

	// Mendapatkan ID dari path parameter
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 404, "Format data pengguna tidak valid")
		return
	}

	// Mencari user berdasarkan ID
	var existingUser m.Users
	result := db.First(&existingUser, userID)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			sendErrorResponse(w, 404, "Pengguna tidak ditemukan")
		} else {
			log.Println(result.Error)
			sendErrorResponse(w, 500, "Kesalahan internal server")
		}
		return
	}

	// Mendapatkan data dari request body
	var updatedUser m.Users
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 400, "Kesalahan saat membaca data update")
		return
	}

	// Memperbarui data user
	result = db.Model(&existingUser).Updates(updatedUser)
	if result.Error != nil {
		log.Println(result.Error)
		sendErrorResponse(w, 500, "Kesalahan saat update data pengguna")
		return
	}

	// Response json
	sendSuccessResponse(w, 200, "Data berhasil diupdate")
}

// 4. DELETE USER (DELETE)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	// Mendapatkan ID dari path parameter
	vars := mux.Vars(r)
	userID := vars["id"]

	// Mengeksekusi pernyataan SQL untuk mendapatkan data user sebelum dihapus
	querySelect := "SELECT Name FROM users WHERE ID=?"
	var userName string
	err := db.QueryRow(querySelect, userID).Scan(&userName)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User dengan ID tersebut tidak ditemukan"))
		return
	}

	// Mengeksekusi pernyataan SQL untuk menghapus user berdasarkan ID
	queryDelete := "DELETE FROM users WHERE ID=?"
	stmt, err := db.Prepare(queryDelete)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	_, err = stmt.Exec(userID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Membuat response JSON
	message := fmt.Sprintf("Data dengan ID %s dan nama %s dihapus", userID, userName)
	var response m.UsersResponse
	response.Status = 200
	response.Message = message

	// Mengirimkan response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 4b. DELETE USER GORM
func DeleteUser2(w http.ResponseWriter, r *http.Request) {
	// Connect gorm dan penangangn error db
	db := gorm_connect(w)

	// Mendapatkan ID dari path parameter menggunakan gorilla/mux
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 400, "invalid user data format")
		return
	}

	// Mencari user berdasarkan ID
	var existingUser m.Users
	result := db.First(&existingUser, userID)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			log.Println(result.Error)
			sendErrorResponse(w, 404, "User not found")
		} else {
			log.Println(result.Error)
			sendErrorResponse(w, 500, "Error fetching user")
		}
		return
	}

	// Menghapus user dari database
	result = db.Delete(&existingUser)
	if result.Error != nil {
		log.Println(result.Error)
		sendErrorResponse(w, 500, "Error deleting user")
		return
	}

	// Respon json
	sendSuccessResponse(w, 200, "Data berhasil dihapus")
}
