package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"go-mysql-api/internal/controller"
	"go-mysql-api/internal/middleware"
	"go-mysql-api/internal/repository"
	"go-mysql-api/internal/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2008-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:123456@tcp(127.0.0.1:3306)/user_db"
	}
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logrus.WithField("error", err).Fatal("Error in opening database connection")
	}

	err = db.Ping()
	if err != nil {
		logrus.WithField("error", err).Fatal("Error in pinging database")
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	fmt.Println("Database connection successful")

	userRepo := &repository.UserRepository{DB: db}
	userServ := &service.UserService{Repo: userRepo}
	userCtrl := &controller.UserController{Service: userServ}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", userCtrl.HandleUsers)
	mux.HandleFunc("/users/", userCtrl.HandleUserDetail)

	handler := middleware.Logging(mux)
	logrus.WithField("port", 8080).Info("Server running at :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logrus.WithField("error", err).Fatal("Error in starting server")
	}
}
