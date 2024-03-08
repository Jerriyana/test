package controllers

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	m "modul4/models"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// func login(w http.ResponseWriter, r *http.Request, db *sql.DB, platform string) error {
//     // Mendapatkan data dari request body
//     var user m.Users
//     err := json.NewDecoder(r.Body).Decode(&user)
//     if err != nil {
//         // Handle error decoding
//         log.Println(err)
//         w.WriteHeader(http.StatusBadRequest)
//         return err
//     }

//     // Memeriksa kredensial user
//     hashedPassword, err := hashPassword(user.Password) // Implementasi hashPassword
//     if err != nil {
//         log.Println(err)
//         w.WriteHeader(http.StatusInternalServerError)
//         return err
//     }

//     var dbUser m.Users
//     err = db.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", user.Email, hashedPassword).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)
//     if err != nil {
//         // Handle error database
//         log.Println(err)
//         if err == sql.ErrNoRows {
//             w.WriteHeader(http.StatusUnauthorized)
//             return err
//         }
//         w.WriteHeader(http.StatusInternalServerError)
//         return err
//     }

//     // Generate token (implementasi generateToken)
//     token, err := generateToken(dbUser.ID)
//     if err != nil {
//         log.Println(err)
//         w.WriteHeader(http.StatusInternalServerError)
//         return err
//     }

//     // Menyiapkan response
//     response := struct {
//         Message string `json:"message"`
//         Token   string `json:"token"`
//     }{
//         Message: fmt.Sprintf("Success login from %s", platform),
//         Token:   token,
//     }

//     // Mengirimkan response
//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusOK)
//     err = json.NewEncoder(w).Encode(response)
//     if err != nil {
//         log.Println(err)
//         return err
//     }

//     return nil
// }
