package controllers

import (
	"database/sql"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connect(w http.ResponseWriter) *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_latihan_pbp?parseTime=true&loc=Asia%2FJakarta")
	if err != nil {
		sendErrorResponse(w, 500, "Internal Server Error! Database Connection Error")
	}
	return db
}

func gorm_connect(w http.ResponseWriter) *gorm.DB {
	dsn := "root:@tcp(localhost:3306)/db_latihan_pbp?parseTime=true&loc=Asia%2FJakarta"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		sendErrorResponse(w, 500, "Internal Server Error! Database Connection Error")
	}
	return db
}
