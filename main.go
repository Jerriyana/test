package main

import (
	"fmt"
	"log"
	"modul5/controllers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// 1. Routing untuk endpoint GET
	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/v2/users", controllers.GetAllUsers2).Methods("GET")        // modul 5
	router.HandleFunc("/v2/users/{age}", controllers.GetUsersByAge).Methods("GET") // modul 5
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/detailTransactions", controllers.GetAllDetailTransactions).Methods("GET")
	router.HandleFunc("/detailTransactions/{id}", controllers.GetSingleDetailTransactions).Methods("GET")

	// 2. Routing untuk endpoint POST
	router.HandleFunc("/users", controllers.InsertUser).Methods("POST")
	router.HandleFunc("/v2/users", controllers.InsertUser2).Methods("POST") // modul 5
	router.HandleFunc("/products", controllers.InsertProduct).Methods("POST")
	router.HandleFunc("/transactions", controllers.InsertTransaction).Methods("POST")
	router.HandleFunc("/transactions2", controllers.InsertTransaction2).Methods("PUT")

	// 3. Routing untuk endpoint PUT
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/v2/users/{id}", controllers.UpdateUser2).Methods("PUT") // modul 5
	router.HandleFunc("/products/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/transactions/{id}", controllers.UpdateTransaction).Methods("PUT")

	// 4. Routing untuk endpoint DELETE
	router.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/v2/users/{id}", controllers.DeleteUser2).Methods("DELETE") // modul 5
	router.HandleFunc("/products/{id}", controllers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/transactions/{id}", controllers.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/singleproducts/{id}", controllers.DeleteSingleProduct).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8890")
	log.Println("Connected to port 8890")
	log.Fatal(http.ListenAndServe(":8890", router))
}
