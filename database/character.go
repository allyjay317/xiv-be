package database

import "github.com/alyjay/xiv-be/types"

type CharacterRow struct {
	LodestoneId string `json:"lodestone_id,omitempty" db:"character_id"`
	UserId      string `json:"user_id,omitempty" db:"user_id"`
	Name        string `json:"name" db,omitempty:"name"`
	Avatar      string `json:"avatar,omitempty" db:"avatar"`
	Portrait    string `json:"portrait,omitempty" db:"portrait"`
}

func InsertCharacter(c CharacterRow) (err error) {
	db, err := GetDb()
	if err != nil {
		return err
	}

	_, err = db.NamedExec(`
	INSERT INTO characters (
		character_id, 
		user_id, 
		name, 
		avatar, 
		portrait
	) VALUES (
	 	:character_id, 
		:user_id, 
		:name, 
		:avatar, 
		:portrait
	)`, c)

	return err
}

func UpdateCharacter(c CharacterRow) (err error) {
	db, err := GetDb()
	if err != nil {
		return err
	}

	_, err = db.NamedExec(`
	UPDATE characters SET
		name = :name, 
		avatar = :avatar, 
		portrait = :portrait
	WHERE character_id = :character_id`, c)

	return err
}

func GetCharacters(id string) (characters []types.Character, err error) {
	db, err := GetDb()
	if err != nil {
		return characters, err
	}
	err = db.Select(&characters, `
		SELECT 
			character_id, 
			name, 
			avatar, 
			portrait 
		FROM characters WHERE user_id=$1`, id)

	return characters, err
}
