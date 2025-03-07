package gearset

import (
	"encoding/json"
	"net/http"

	database "github.com/alyjay/xiv-be/database"
	"github.com/alyjay/xiv-be/response"
	types "github.com/alyjay/xiv-be/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func SetUpGearSetRoutes(r *mux.Router) {
	r.HandleFunc("/{characterId}", AddGearSet).Methods("POST")
	r.HandleFunc("/{characterId}", GetGearSets).Methods("GET")
	r.HandleFunc("/{characterId}", BulkUpdateGearSets).Methods("PATCH")
	r.HandleFunc("/{characterId}/{id}", UpdateGearSet).Methods("PATCH")
	r.HandleFunc("/{characterId}/{id}", DeleteGearSet).Methods("DELETE")
}

func AddGearSet(w http.ResponseWriter, r *http.Request) {
	var req types.GearSetV2
	_ = json.NewDecoder(r.Body).Decode(&req)

	characterId := mux.Vars(r)["characterId"]

	newUUID := uuid.NewString()

	set, err := database.InsertGearSetV2(types.GearSetV2{
		ID:          newUUID,
		UserId:      req.UserId,
		CharacterId: characterId,
		Name:        req.Name,
		Job:         req.Job,
		Items:       req.Items,
		Index:       req.Index,
	})
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	set.ID = newUUID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(set)
}

func GetGearSets(w http.ResponseWriter, r *http.Request) {
	archived := r.URL.Query().Has("archived")

	characterId := mux.Vars(r)["characterId"]
	var Sets []types.GearSetV2

	Sets, err := database.SelectGearSetsForCharacterV2(characterId, archived)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Sets)
}

func UpdateGearSet(w http.ResponseWriter, r *http.Request) {

	var req types.GearSetV2
	_ = json.NewDecoder(r.Body).Decode(&req)

	id := mux.Vars(r)["id"]
	characterId := mux.Vars(r)["characterId"]
	req.CharacterId = characterId
	req.ID = id

	err := database.UpdateGearSetV2(req)

	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
}

func DeleteGearSet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := database.DeleteGearSet(id)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func BulkUpdateGearSets(w http.ResponseWriter, r *http.Request) {
	characterId := mux.Vars(r)["characterId"]

	var req []types.GearSetV2
	_ = json.NewDecoder(r.Body).Decode(&req)

	for i, s := range req {
		database.UpdateGearSetV2(types.GearSetV2{
			ID:          s.ID,
			Name:        s.Name,
			Job:         s.Job,
			Index:       i,
			Items:       s.Items,
			CharacterId: characterId,
			Archived:    s.Archived,
		})
	}

}
