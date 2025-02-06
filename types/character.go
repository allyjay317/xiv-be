package types

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
