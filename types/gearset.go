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

type Items map[int]interface{}

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
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Job   Job    `json:"job" db:"job"`
	Items Items  `json:"items" db:"config"`
}

type GearSetRequest struct {
	UserId string `json:"id" db:"user_id"`
	Name   string `json:"name" db:"name"`
	Job    Job    `json:"job" db:"job"`
	Tier   string `json:"tier" db:"tier"`
	Items  map[int]struct {
		Augmented bool   `json:"augmented"`
		Have      bool   `json:"have"`
		Source    Source `json:"source"`
	} `json:"items" db:"config"`
}

type AddGearSetResponse struct {
	ID string `json:"id"`
}
