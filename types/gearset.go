package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Source int
type Job int
type Slot int

const (
	Raid Source = iota
	Crafted
	Tome
	Chaotic
	Ultimate
)

const (
	WHM Job = iota
	SCH
	AST
	SGE
	PLD
	WAR
	DRK
	GNB
	MNK
	DRG
	NIN
	SAM
	RPR
	BRD
	MCH
	DNC
	RDM
	BLM
	SMN
	BLU
	VPR
	PCT
)

const (
	HEAD Slot = iota
	BODY
	HANDS
	LEGS
	FEET
	EARRINGS
	NECKLACE
	BRACELET
	RING1
	RING2
	WEAPON
	OFFHAND
)

type GearPiece struct {
	Id        int    `json:"id" db:"id"`
	Source    Source `json:"source" db:"source"`
	Have      bool   `json:"have" db:"have"`
	Augmented bool   `json:"augmented,omitempty" db:"augmented"`
	Priority  int    `json:"priority" db:"priority"`
}

type GearPieceRow struct {
	GearPiece
	Slot Slot `json:"slot" db:"slot"`
}

type Items map[int]interface {
}

func (i *Items) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &i)
		return nil
	case string:
		json.Unmarshal([]byte(v), &i)
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type %T", v)
	}
}

func (i *Items) Value() (driver.Value, error) {
	return json.Marshal(i)
}

type ItemsV2 = map[Slot]GearPiece

type GearSet struct {
	Archived    bool   `json:"archived" db:"archived"`
	ID          string `json:"id" db:"id"`
	Index       int    `json:"index" db:"index"`
	Items       Items  `json:"items" db:"config"`
	Job         Job    `json:"job" db:"job"`
	Name        string `json:"name" db:"name"`
	UserId      string `json:"user_id" db:"user_id"`
	CharacterId string `json:"character_id" db:"character_id"`
}

type GearSetV2 struct {
	Archived    bool    `json:"archived" db:"archived"`
	ID          string  `json:"id" db:"id"`
	Index       int     `json:"index" db:"index"`
	Job         Job     `json:"job" db:"job"`
	Name        string  `json:"name" db:"name"`
	UserId      string  `json:"user_id" db:"user_id"`
	CharacterId string  `json:"character_id" db:"character_id"`
	Items       ItemsV2 `json:"items" db:"config"`
}

type GearSetRow struct {
	GearSet
	Items []byte `json:"id" db:"config"`
}

type AddGearSetResponse struct {
	ID string `json:"id"`
}
