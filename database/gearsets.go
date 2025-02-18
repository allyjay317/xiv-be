package database

import (
	"encoding/json"

	"github.com/alyjay/xiv-be/types"
)

func InsertGearSet(g types.GearSet) (err error) {
	db, err := GetDb()

	if err != nil {
		return err
	}

	items, err := json.Marshal(g.Items)
	if err != nil {
		return err
	}

	gs := types.GearSetRow{
		GearSet: g,
		Items:   items,
	}

	_, err = db.NamedExec(`
	INSERT INTO gear_sets (
		id,
		user_id,
		character_id,
		name,
		job,
		config,
		index
	) VALUES (
		 :id,
		 :user_id,
		 :character_id,
		 :name,
		 :job,
		 :config,
		 :index
	)
	`, gs)

	return err
}

func SelectGearSetsForCharacter(id string, archived bool) (gs []types.GearSet, err error) {
	db, err := GetDb()

	if err != nil {
		return gs, err
	}

	err = db.Select(&gs, `
	SELECT 
		id, name, job, config
	FROM gear_sets WHERE
		character_id = $1
	AND
		archived = $2
	ORDER BY index ASC
	`, id, archived)

	return gs, err
}

func UpdateGearSet(g ...types.GearSet) (err error) {
	db, err := GetDb()

	if err != nil {
		return err
	}
	var gs []types.GearSetRow

	for _, s := range g {
		items, err := json.Marshal(s.Items)
		if err != nil {
			return err
		}
		gs = append(gs, types.GearSetRow{
			GearSet: s,
			Items:   items,
		})
	}

	_, err = db.NamedExec(`
		UPDATE gear_sets SET
			name = :name,
			job = :job,
			config = :config,
			index = :index,
			archived = :archived
		WHERE id = :id AND character_id = :character_id
	`, gs)

	return err
}

func DeleteGearSet(id string) (err error) {
	db, err := GetDb()

	if err != nil {
		return err
	}

	_, err = db.Exec(`
	DELETE FROM gear_sets 
		WHERE id = $1
	`, id)

	return err
}
