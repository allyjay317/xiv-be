package types

import (
	"errors"
)

type PurchaseInfo struct {
	Tome int
	Raid int
}

var one = PurchaseInfo{
	Tome: 375,
	Raid: 3,
}

var two = PurchaseInfo{
	Tome: 495,
	Raid: 4,
}

var three = PurchaseInfo{
	Tome: 825,
	Raid: 8,
}

var four = PurchaseInfo{
	Tome: 500,
	Raid: 8,
}

var slotInfo = map[Slot]PurchaseInfo{
	EARRINGS: one,
	NECKLACE: one,
	BRACELET: one,
	RING1:    one,
	RING2:    one,
	HEAD:     two,
	HANDS:    two,
	FEET:     two,
	BODY:     three,
	LEGS:     three,
	WEAPON:   four,
}

func GetCost(slot Slot, source Source) (int, error) {
	var info = slotInfo[slot]
	if source == Raid {
		return info.Raid, nil
	}
	if source == Tome {
		return info.Tome, nil
	}
	return -1, errors.New("not a valid purchase")
}

type PurchaseRequest struct {
	CharacterId string
	Slot        Slot
	Source      Source
	Date        string
	Id          int
}
