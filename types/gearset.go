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
	Source    Source `json:"source"`
	Have      bool   `json:"have"`
	Augmented bool   `json:"augmented,omitempty"`
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
	default:
		return fmt.Errorf("unsupported type %T", v)
	}
}

func (i *Items) Value() (driver.Value, error) {
	return json.Marshal(i)
}

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

type GearSetRow struct {
	GearSet
	Items []byte `json:"id" db:"config"`
}

type AddGearSetResponse struct {
	ID string `json:"id"`
}
