package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alyjay/xiv-be/character"
	"github.com/alyjay/xiv-be/database"
	"github.com/alyjay/xiv-be/gearset"
	"github.com/alyjay/xiv-be/user"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	database.Migrate()

	router := mux.NewRouter()
	router.HandleFunc("/user", user.GetUser).Methods("GET")
	router.HandleFunc("/login", user.LoginUser).Methods("GET")
	rCharacter := router.PathPrefix("/character").Subrouter()
	character.SetUpCharacterRoutes(rCharacter)
	rGearset := router.PathPrefix("/gearset").Subrouter()
	gearset.SetUpGearSetRoutes(rGearset)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	})

	log.Fatal(http.ListenAndServe(":"+port, c.Handler(router)))
}
