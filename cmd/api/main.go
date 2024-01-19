package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/mrbaloch555/go-chi-auth/models"
)

var webPort = "8080"

type Config struct {
	DB     *sql.DB
	Models *models.Models
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DSN := os.Getenv("DSN")
	conn, err := connectDB(DSN)
	if err != nil {
		log.Panic(conn)
	}
	app := Config{
		DB:     conn,
		Models: models.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	fmt.Printf("Serving app on port %s\n", webPort)
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func connectDB(DSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", DSN)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres")

	return db, nil
}
