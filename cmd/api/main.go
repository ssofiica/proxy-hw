package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssofiica/proxy-hw/internal/api"
	"github.com/ssofiica/proxy-hw/internal/repo"
)

var (
	API_PORT      = ":8000"
	POSTGRES_CONN = "postgres://svalova:mydbpass@localhost:5432/test-gaz"
)

func main() {
	db, err := pgxpool.New(context.Background(), POSTGRES_CONN)
	if err != nil {
		fmt.Println("error wih db", err)
	}
	repository := repo.NewRepo(db)
	handler := api.NewHandler(repository)

	r := mux.NewRouter()
	r.HandleFunc("/requests", handler.GetAll).Methods("GET")
	r.HandleFunc("/requests/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/repeat/{id}", handler.Repeat).Methods("GET")
	r.HandleFunc("/scan/{id}", handler.Scan).Methods("GET")

	apiServer := &http.Server{
		Addr:    API_PORT,
		Handler: r,
	}
	if err := apiServer.ListenAndServe(); err != nil {
		log.Fatal("proxyServer wasn't started:", err)
	}

}
