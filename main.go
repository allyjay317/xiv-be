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
	log.Println("Running Migrations")
	err := database.Migrate()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Migrations Completed")

	log.Println("Setting Up Routes")
	router := mux.NewRouter()
	router.HandleFunc("/user", user.GetUser).Methods("GET")
	// router.HandleFunc("/users", user.CreateUser).Methods("POST")
	router.HandleFunc("/login", user.LoginUser).Methods("GET")
	router.HandleFunc("/character", character.SearchCharacter).Methods("POST")
	router.HandleFunc("/character/verify", character.VerifyCharacter).Methods("POST")
	router.HandleFunc("/gearset/{characterId}", gearset.AddGearSet).Methods("POST")
	router.HandleFunc("/gearset/{characterId}/{id}", gearset.UpdateGearSet).Methods("PATCH")
	router.HandleFunc("/gearset/{characterId}/{id}", gearset.DeleteGearSet).Methods("DELETE")
	router.HandleFunc("/gearset/{characterId}", gearset.GetGearSets).Methods("GET")

	log.Println("Setting up CORS")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	})

	log.Fatal(http.ListenAndServe(":"+port, c.Handler(router)))
}
