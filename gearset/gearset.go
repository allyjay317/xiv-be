package gearset

import (
	"encoding/json"
	"net/http"

	database "github.com/alyjay/xiv-be/database"
	types "github.com/alyjay/xiv-be/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func SetUpGearSetRoutes(r *mux.Router) {
	r.HandleFunc("/{characterId}", AddGearSet).Methods("POST")
	r.HandleFunc("/{characterId}/{id}", UpdateGearSet).Methods("PATCH")
	r.HandleFunc("/{characterId}/{id}", DeleteGearSet).Methods("DELETE")
	r.HandleFunc("/{characterId}", GetGearSets).Methods("GET")
}

func AddGearSet(w http.ResponseWriter, r *http.Request) {
	var req types.GearSetRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	characterId := mux.Vars(r)["characterId"]

	db, err := database.GetDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue" + err.Error()))
		return
	}

	newUUID := uuid.NewString()

	items, err := json.Marshal(req.Items)

	_, err = db.Exec(`
		INSERT INTO gear_sets (
			id, 
			user_id, 
			character_id, 
			name, 
			job, 
			config
		) VALUES (
		 	$1, 
			$2, 
			$3, 
			$4, 
			$5,
			$6)`,
		newUUID,
		req.UserId,
		characterId,
		req.Name,
		req.Job,
		items)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error entering into database"))
	}

	var res types.AddGearSetResponse
	res.ID = newUUID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func GetGearSets(w http.ResponseWriter, r *http.Request) {
	characterId := mux.Vars(r)["characterId"]
	var Sets []types.GearSet

	db, err := database.GetDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue" + err.Error()))
		return
	}

	err = db.Select(&Sets, `SELECT id, name, job, config from gear_sets WHERE character_id = $1 AND archived = false`, characterId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Sets)
}

func UpdateGearSet(w http.ResponseWriter, r *http.Request) {
	var req types.GearSetRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	id := mux.Vars(r)["id"]
	db, err := database.GetDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue" + err.Error()))
		return
	}

	items, err := json.Marshal(req.Items)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	_, err = db.Exec(`UPDATE gear_sets SET name = $1, job = $2, config = $3 WHERE id = $4`,
		req.Name,
		req.Job,
		items,
		id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not update"))
		return
	}
}

func DeleteGearSet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	db, err := database.GetDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue" + err.Error()))
		return
	}

	_, err = db.Exec(`DELETE FROM gear_sets WHERE id = $1`, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
