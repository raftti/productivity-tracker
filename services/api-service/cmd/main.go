package main

import (
	"fmt"
	"log"
	"net/http"

	dbClient "api-service/internal/client/db"
	userHandler "api-service/internal/handlers/user"

	"github.com/go-chi/chi/v5"
)

func main() {
	dbClient, dbConn := dbClient.NewDBServiceClient()
	defer dbConn.Close()
	
	r := chi.NewRouter()

	userRoutes := userHandler.NewUserHandler(dbClient)

	userRoutes.RegisterRoutes(r)

	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}