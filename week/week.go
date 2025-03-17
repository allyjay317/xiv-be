package week

import (
	"encoding/json"
	"net/http"

	"github.com/alyjay/xiv-be/database"
	"github.com/alyjay/xiv-be/response"
	"github.com/alyjay/xiv-be/types"
	"github.com/gorilla/mux"
)

func SetUpWeekRoutes(r *mux.Router) {
	r.HandleFunc("/{id}", MakePurchase).Methods("POST")
}

func MakePurchase(w http.ResponseWriter, r *http.Request) {
	characterId := mux.Vars(r)["id"]
	var req types.PurchaseRequest
	json.NewDecoder(r.Body).Decode(&req)
	req.CharacterId = characterId

	database.InsertPurchase(req)
	piece, err := database.GetGearPiece(req.Id)
	if err != nil {
		response.BadRequestError(w, err.Error())
		return
	}
	piece.Have = true
	database.UpdateGearPiece(piece)

	cost, err := types.GetCost(req.Slot, req.Source)
	if err != nil {
		response.BadRequestError(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(cost)
}

func GetPurchases(w http.ResponseWriter, r *http.Request) {

}
