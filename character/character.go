package character

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	database "github.com/alyjay/xiv-be/database"
	"github.com/gorilla/mux"
	"github.com/karashiiro/bingode"
	"github.com/xivapi/godestone/v2"
)

type AddCharacterRequest struct {
	LodestoneId string `json:"lodestone_id"`
}

type VerifyCharacterRequest struct {
	ID          string `json:"id"`
	LodestoneId string `json:"lodestone_id"`
	VerifyCode  string `json:"verify_code"`
}

type Character struct {
	CharacterId string `json:"id" db:"character_id"`
	Name        string `json:"name" db:"name"`
	Avatar      string `json:"avatar" db:"avatar"`
	Portrait    string `json:"portrait" db:"portrait"`
}

func SetUpCharacterRoutes(r *mux.Router) {
	r.HandleFunc("", SearchCharacter).Methods("POST")
	r.HandleFunc("/verify", VerifyCharacter).Methods("POST")
	r.HandleFunc("/c/{id}", DeleteCharacter).Methods("DELETE")
	r.HandleFunc("/c/{id}", UpdateCharacter).Methods("GET")
}

func FetchCharacterData(id uint32) (c *godestone.Character, err error) {
	s := godestone.NewScraper(bingode.New(), godestone.EN)
	c, err = s.FetchCharacter(id)
	return c, err
}

func SearchCharacter(w http.ResponseWriter, r *http.Request) {
	var req AddCharacterRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	id, err := strconv.ParseUint(req.LodestoneId, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error"))
		return
	}

	c, err := FetchCharacterData(uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}

	var chara Character
	chara.Avatar = c.Avatar
	chara.CharacterId = fmt.Sprintf("%d", c.ID)
	chara.Name = c.Name
	chara.Portrait = c.Portrait

	json.NewEncoder(w).Encode(chara)
}

func VerifyCharacter(w http.ResponseWriter, r *http.Request) {
	var req VerifyCharacterRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	id, err := strconv.ParseUint(req.LodestoneId, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error"))
		return
	}

	c, err := FetchCharacterData(uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}

	if c.Bio != req.VerifyCode {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Verify String Incorrect"))
		return
	}

	db, err := database.GetDb(w)
	if err != nil {
		return
	}

	_, err = db.Exec(`INSERT INTO characters (character_id, user_id, name, avatar, portrait) VALUES ($1, $2, $3, $4, $5)`, req.LodestoneId, req.ID, c.Name, c.Avatar, c.Portrait)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Verified"))
}

func DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	db, err := database.GetDb(w)
	if err != nil {
		return
	}

	_, err = db.Exec(`DELETE FROM characters WHERE character_id = $1`, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	characterId := mux.Vars(r)["id"]

	id, err := strconv.ParseUint(characterId, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error"))
		return
	}

	c, err := FetchCharacterData(uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}

	db, err := database.GetDb(w)
	if err != nil {
		return
	}

	_, err = db.Exec(`UPDATE characters SET
	name = $1, 
	avatar = $2, 
	portrait = $3
	WHERE character_id = $4`,
		c.Name,
		c.Avatar,
		c.Portrait,
		id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Character{
		Name:        c.Name,
		Avatar:      c.Avatar,
		Portrait:    c.Portrait,
		CharacterId: fmt.Sprint(c.ID),
	})
}
