package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	m "modul5/models"
	"net/http"

	"github.com/gorilla/mux"
)

// 1. GET ALL PRODUCTS
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	query := "SELECT * FROM products"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "select error")
		return
	}

	var product m.Products
	var products []m.Products
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Println(err)
			sendErrorResponse(w, 500, "db error")
			return
		} else {
			products = append(products, product)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.ProductsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = products // products
	json.NewEncoder(w).Encode(response)
}

// 2. INSERT PRODUCT
func InsertProduct(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	query := "INSERT INTO products (ID, Name, Price) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(query)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "insert error")
		return
	}

	// Mendapatkan data dari request body
	var product m.Products
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		sendErrorResponse(w, 500, "get data error")
		log.Println(err)
		return
	}

	// Menjalankan pernyataan Exec dengan data dari request body
	_, err = stmt.Exec(product.ID, product.Name, product.Price)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, 500, "eksekusi error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User berhasil dibuat"))

	var products []m.Products
	products = append(products, product)

	var response m.ProductsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = products
	json.NewEncoder(w).Encode(response)
}

// 3. UPDATE PRODUCT
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	// Mendapatkan ID dari path parameter
	vars := mux.Vars(r)
	productID := vars["id"]

	// Membaca data dari request body
	var updatedProduct m.Products
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		sendErrorResponse(w, 500, "gagal membaca data")
		log.Println(err)
		return
	}

	// Mengeksekusi pernyataan SQL untuk mengupdate product berdasarkan ID
	query := "UPDATE products SET Name=?, Price=? WHERE ID=?"
	stmt, err := db.Prepare(query)
	if err != nil {
		sendErrorResponse(w, 500, "gagal update product berdasarkan id")
		log.Println(err)
		return
	}

	_, err = stmt.Exec(updatedProduct.Name, updatedProduct.Price, productID)
	if err != nil {
		sendErrorResponse(w, 500, "gagal eksekusi")
		log.Println(err)
		return
	}

	// Mengeksekusi pernyataan SQL untuk mendapatkan data product yang telah diupdate
	querySelect := "SELECT ID, Name, Price FROM products WHERE ID=?"
	rows, err := db.Query(querySelect, productID)
	if err != nil {
		sendErrorResponse(w, 500, "gagal mendapatkan data yang telah diupdate")
		log.Println(err)
		return
	}
	defer rows.Close()

	// Mengonversi hasil query ke dalam bentuk array products
	var products []m.Products
	for rows.Next() {
		var product m.Products
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			sendErrorResponse(w, 500, "gagal konversi ke array")
			log.Println(err)
			return
		}
		products = append(products, product)
	}

	// Membuat response JSON
	var response m.ProductsResponse
	response.Status = 200
	response.Message = "Update success"
	response.Data = products
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 4. DELETE PRODUCT
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	// Mendapatkan ID dari path parameter
	vars := mux.Vars(r)
	productID := vars["id"]

	// Mengeksekusi pernyataan SQL untuk mendapatkan data product sebelum dihapus
	querySelect := "SELECT Name FROM products WHERE ID=?"
	var productName string
	err := db.QueryRow(querySelect, productID).Scan(&productName)
	if err != nil {
		// send error response
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User dengan ID tersebut tidak ditemukan"))
		return
	}

	// Mengeksekusi pernyataan SQL untuk menghapus product berdasarkan ID
	queryDelete := "DELETE FROM products WHERE ID=?"
	stmt, err := db.Prepare(queryDelete)
	if err != nil {
		sendErrorResponse(w, 500, "gagal hapus product berdasarkan ID")
		log.Println(err)
		return
	}

	_, err = stmt.Exec(productID)
	if err != nil {
		sendErrorResponse(w, 500, "gagal eksekusi")
		log.Println(err)
		return
	}

	// Membuat response JSON
	message := fmt.Sprintf("Data dengan ID %s dan nama %s dihapus", productID, productName)
	var response m.ProductsResponse
	response.Status = 200
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 4b. DELETE PRODUCT AND PRODUCT IN TRANSACTIONS
func DeleteSingleProduct(w http.ResponseWriter, r *http.Request) {
	db := connect(w)
	defer db.Close()

	// Mendapatkan ID dari path parameter
	vars := mux.Vars(r)
	productID := vars["id"]

	// Delete transactions related to the product
	tx, err := db.Begin()
	if err != nil {
		sendErrorResponse(w, 500, "internal error")
		return
	}

	stmt, err := tx.Prepare("DELETE FROM transactions WHERE ProductID = ?")
	if err != nil {
		tx.Rollback() // Rollback transaction if error occurs
		sendErrorResponse(w, 500, "gagal delete")
		return
	}

	_, err = stmt.Exec(productID)
	if err != nil {
		tx.Rollback() // Rollback transaction if error occurs
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Delete product
	stmt, err = tx.Prepare("DELETE FROM products WHERE ID = ?")
	if err != nil {
		tx.Rollback() // Rollback transaction if error occurs
		sendErrorResponse(w, 500, "gagal delete berdasarkan id")
		return
	}

	_, err = stmt.Exec(productID)
	if err != nil {
		tx.Rollback() // Rollback transaction if error occurs
		sendErrorResponse(w, 500, "internal error")
		return
	}

	err = tx.Commit() // Commit transaction if all operations succeed
	if err != nil {
		sendErrorResponse(w, 500, "internal error")
		return
	}

	// Membuat response JSON
	var response m.ProductsResponse
	response.Status = 200
	response.Message = "delete succes"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
