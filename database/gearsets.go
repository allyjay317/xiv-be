package database

import "github.com/alyjay/xiv-be/types"

func InsertGearSet(g types.GearSet) (err error) {
	db, err := GetDb()

	if err != nil {
		return err
	}

	_, err = db.NamedExec(`
	INSERT INTO gear_sets (
		id,
		user_id,
		character_id,
		name,
		job,
		config
	)
	`, g)

	return err
}

func SelectGearSet(id string) (err error) {
	db, err := GetDb()

	if err != nil {
		return err
	}

	_, err = db.Exec(`
	SELECT 
		id, name, job, config
	FROM gear_sets WHERE
		character_id = $1
	AND
		archived = false
	`, id)

	return err
}

func UpdateGearSet(g types.GearSet) (err error) {
	db, err := GetDb()

	if err != nil {
		return err
	}

	_, err = db.NamedExec(`
		UPDATE gear_sets SET
			name = :name
			job = :job
			config = :config
			index = :index
		WHERE id = :id
	`, g)

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
