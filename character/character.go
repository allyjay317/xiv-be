package character

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	database "github.com/alyjay/xiv-be/database"
	"github.com/alyjay/xiv-be/response"
	types "github.com/alyjay/xiv-be/types"
	"github.com/gorilla/mux"
	"github.com/karashiiro/bingode"
	"github.com/xivapi/godestone/v2"
)

func SetUpCharacterRoutes(r *mux.Router) {
	r.HandleFunc("", SearchCharacter).Methods("POST")
	r.HandleFunc("/verify", VerifyCharacter).Methods("POST")
	r.HandleFunc("/c/{id}", DeleteCharacter).Methods("DELETE")
	r.HandleFunc("/c/{id}", UpdateCharacter).Methods("GET")
}

func FetchCharacterData(strId string) (c *godestone.Character, err error) {
	id, err := strconv.ParseUint(strId, 10, 32)
	if err != nil {
		return nil, err
	}
	s := godestone.NewScraper(bingode.New(), godestone.EN)
	c, err = s.FetchCharacter(uint32(id))
	return c, err
}

func SearchCharacter(w http.ResponseWriter, r *http.Request) {
	var req types.AddCharacterRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	c, err := FetchCharacterData(req.LodestoneId)
	if err != nil {
		response.NotFoundError(w, "Failed to get character data")
		return
	}

	var chara types.Character
	chara.Avatar = c.Avatar
	chara.CharacterId = fmt.Sprintf("%d", c.ID)
	chara.Name = c.Name
	chara.Portrait = c.Portrait

	json.NewEncoder(w).Encode(chara)
}

func VerifyCharacter(w http.ResponseWriter, r *http.Request) {
	var req types.VerifyCharacterRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	c, err := FetchCharacterData(req.LodestoneId)
	if err != nil {
		response.NotFoundError(w, "Failed to get character data")
		return
	}

	if c.Bio != req.VerifyCode {
		response.NotFoundError(w, "Verify String Incorrect")
		return
	}

	err = database.InsertCharacter(database.CharacterRow{
		LodestoneId: req.LodestoneId,
		UserId:      req.ID,
		Name:        c.Name,
		Avatar:      c.Avatar,
		Portrait:    c.Portrait,
	})

	if err != nil {
		response.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Verified"))
}

func DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	db, err := database.GetDb()
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	_, err = db.Exec(`DELETE FROM characters WHERE character_id = $1`, id)
	if err != nil {
		response.BadRequestError(w)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	characterId := mux.Vars(r)["id"]

	c, err := FetchCharacterData(characterId)
	if err != nil {
		response.NotFoundError(w, "Failed to retrieve character data")
		return
	}

	err = database.UpdateCharacter(database.CharacterRow{
		LodestoneId: characterId,
		Name:        c.Name,
		Avatar:      c.Avatar,
		Portrait:    c.Portrait,
	})
	if err != nil {
		response.InternalServerError(w, "Failed to update character")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Character{
		Name:        c.Name,
		Avatar:      c.Avatar,
		Portrait:    c.Portrait,
		CharacterId: fmt.Sprint(c.ID),
	})
}
