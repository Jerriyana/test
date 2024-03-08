package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	m "modul5/models"
	"net/http"

	"github.com/gorilla/mux"
)

// 1. GET ALL TRANSACTIONS
func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	query := "SELECT * FROM transactions"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	var transaction m.Transactions
	var transactions []m.Transactions
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.ProductID, &transaction.Quantity, &transaction.UserID); err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "internal error")
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

// 1b. GET ALL DETAIL TRANSACTIONS
func GetAllDetailTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	query := `
	SELECT
	  u.ID, u.Name, u.Age, u.Address,
	  t.ID, t.Quantity,
	  p.ID, p.Name, p.Price
	FROM users u
	JOIN transactions t ON u.ID = t.UserID
	JOIN products p ON t.ProductID = p.ID
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	var detailTransaction m.DetailTransactions
	var detailTransactions []m.DetailTransactions
	for rows.Next() {
		if err := rows.Scan(&detailTransaction.Users.ID, &detailTransaction.Users.Name,
			&detailTransaction.Users.Age, &detailTransaction.Users.Address,
			&detailTransaction.Transactions.ID, &detailTransaction.Transactions.Quantity,
			&detailTransaction.Products.ID, &detailTransaction.Products.Name, &detailTransaction.Products.Price); err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "internal error")
			return
		} else {
			detailTransactions = append(detailTransactions, detailTransaction)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.DetailTransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = detailTransactions
	json.NewEncoder(w).Encode(response)
}

func GetSingleDetailTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	vars := mux.Vars(r)
	userID := vars["id"]

	query := `
	SELECT
	  u.ID, u.Name, u.Age, u.Address,
	  t.ID, t.Quantity,
	  p.ID, p.Name, p.Price
	FROM users u
	JOIN transactions t ON u.ID = t.UserID
	JOIN products p ON t.ProductID = p.ID
	WHERE u.ID=?
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	var detailTransaction m.DetailTransactions
	var detailTransactions []m.DetailTransactions
	for rows.Next() {
		if err := rows.Scan(&detailTransaction.Users.ID, &detailTransaction.Users.Name,
			&detailTransaction.Users.Age, &detailTransaction.Users.Address,
			&detailTransaction.Transactions.ID, &detailTransaction.Transactions.Quantity,
			&detailTransaction.Products.ID, &detailTransaction.Products.Name, &detailTransaction.Products.Price); err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "internal error")
			return
		} else {
			detailTransactions = append(detailTransactions, detailTransaction)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.DetailTransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = detailTransactions
	json.NewEncoder(w).Encode(response)
}

// 2. INSERT TRANSACTION (POST)
func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	query := "INSERT INTO transactions (ID, UserID, ProductID, Quantity) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Mendapatkan data dari request body
	var transaction m.Transactions
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Menjalankan pernyataan Exec dengan data dari request body
	_, err = stmt.Exec(transaction.ID, transaction.UserID, transaction.ProductID, transaction.Quantity)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User berhasil dibuat"))

	// Mengisi nilai objek transactions
	var transactions []m.Transactions
	transactions = append(transactions, transaction)

	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions

	// Mengirimkan response JSON
	json.NewEncoder(w).Encode(response)
}

// INSERT TRANSACTION 2 (POST)
func InsertTransaction2(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	// Mendapatkan data dari request body
	var transaction m.Transactions
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Memeriksa apakah produk sudah ada
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE ID = ?)", transaction.ProductID).Scan(&exists)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return // Mengirimkan response error
	}

	if !exists {
		// Memasukan produk baru dengan hanya id
		stmt, err := db.Prepare("INSERT INTO products (ID) VALUES(?)")
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "internal error")
			return
		}

		_, err = stmt.Exec(transaction.ProductID)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "internal error")
			return
		}
	}

	// insert ke transactions
	query := "INSERT INTO transactions (ID, UserID, ProductID, Quantity) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Menjalankan pernyataan Exec dengan data dari request body
	_, err = stmt.Exec(transaction.ID, transaction.UserID, transaction.ProductID, transaction.Quantity)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User berhasil dibuat"))

	// Mengisi nilai objek transactions
	var transactions []m.Transactions
	transactions = append(transactions, transaction)

	// Menentukan pesan sesuai keberadaan produk
	var message string
	if exists {
		message = "transaksi berhasil dibuat"
	} else {
		message = "transaksi berhasil dibuat, meski sebelumnya tidak ada product tersebut di dalam tabel products"
	}

	var response m.TransactionsResponse
	response.Status = 200
	response.Message = message
	response.Data = transactions

	// Mengirimkan response JSON
	json.NewEncoder(w).Encode(response)
}

// 3. UPDATE TRANSACTION (PUT)
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	// Mendapatkan ID dari path parameter
	vars := mux.Vars(r)
	transactionID := vars["id"]

	// Membaca data dari request body
	var updatedTransaction m.Transactions
	err := json.NewDecoder(r.Body).Decode(&updatedTransaction)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Mengeksekusi pernyataan SQL untuk mengupdate transactions berdasarkan ID
	query := "UPDATE transactions SET UserID=?, ProductID=?, Quantity=? WHERE ID=?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	_, err = stmt.Exec(updatedTransaction.UserID, updatedTransaction.ProductID, updatedTransaction.Quantity, transactionID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Mengeksekusi pernyataan SQL untuk mendapatkan data transaction yang telah diupdate
	querySelect := "SELECT ID, UserID, ProductID, Quantity FROM transactions WHERE ID=?"
	rows, err := db.Query(querySelect, transactionID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}
	defer rows.Close()

	// Mengonversi hasil query ke dalam bentuk array transaction
	var transactions []m.Transactions
	for rows.Next() {
		var transaction m.Transactions
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "internal error")
			return
		}
		transactions = append(transactions, transaction)
	}

	// Membuat response JSON
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Update success"
	response.Data = transactions

	// Mengirimkan response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 4. DELETE TRANSACTION (DELETE)
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	// Mendapatkan ID dari path parameter
	vars := mux.Vars(r)
	transactionID := vars["id"]

	// Mengeksekusi pernyataan SQL untuk menghapus transaction berdasarkan ID
	queryDelete := "DELETE FROM transactions WHERE ID=?"
	stmt, err := db.Prepare(queryDelete)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	_, err = stmt.Exec(transactionID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Membuat response JSON
	message := fmt.Sprintf("Data dengan ID %s telah dihapus", transactionID)
	var response m.TransactionResponse
	response.Status = 200
	response.Message = message

	// Mengirimkan response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
