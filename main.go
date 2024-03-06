package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/javito2003/crud-go/handlers"
	"github.com/javito2003/crud-go/middlewares"
	"github.com/javito2003/crud-go/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error to load .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JwtSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal("Error to init server:", err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	authRouter := r.NewRoute().Subrouter()
	userRoter := r.NewRoute().Subrouter()
	postRouter := r.NewRoute().Subrouter()

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	authRouter.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)

	userRoter.Use(middlewares.CheckAuthMiddleware(s))
	userRoter.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)

	postRouter.Use(middlewares.CheckAuthMiddleware(s))
	postRouter.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
}
