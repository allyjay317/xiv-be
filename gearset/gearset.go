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
	r.HandleFunc("/{characterId}/{id}", UpdateGearSet).Methods("PATCH")
	r.HandleFunc("/{characterId}/{id}", DeleteGearSet).Methods("DELETE")
	r.HandleFunc("/{characterId}", GetGearSets).Methods("GET")
	r.HandleFunc("/{characterId}", BulkUpdateGearSets).Methods("PATCH")
}

func AddGearSet(w http.ResponseWriter, r *http.Request) {
	var req types.GearSet
	_ = json.NewDecoder(r.Body).Decode(&req)

	characterId := mux.Vars(r)["characterId"]

	newUUID := uuid.NewString()

	err := database.InsertGearSet(types.GearSet{
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

	var res types.AddGearSetResponse
	res.ID = newUUID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func GetGearSets(w http.ResponseWriter, r *http.Request) {
	characterId := mux.Vars(r)["characterId"]
	var Sets []types.GearSet

	Sets, err := database.SelectGearSetsForCharacter(characterId)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Sets)
}

func UpdateGearSet(w http.ResponseWriter, r *http.Request) {

	var req types.GearSet
	_ = json.NewDecoder(r.Body).Decode(&req)

	id := mux.Vars(r)["id"]
	characterId := mux.Vars(r)["characterId"]
	req.CharacterId = characterId
	req.ID = id

	err := database.UpdateGearSet(req)

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

	var req []types.GearSet
	_ = json.NewDecoder(r.Body).Decode(&req)

	for i, s := range req {
		// items, _ := json.Marshal(s.Items)
		database.UpdateGearSet(types.GearSet{
			ID:          s.ID,
			Name:        s.Name,
			Job:         s.Job,
			Index:       i,
			Items:       s.Items,
			CharacterId: characterId,
		})
	}

}
