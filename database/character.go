package database

var insertString = `
	INSERT INTO characters 
		(character_id, 
		user_id, 
		name, 
		avatar, 
		portrait) 
		VALUES (:character_id, :user_id, :name, :avatar, :portrait)`

var updateString = `UPDATE characters SET
	name = :name, 
	avatar = :avatar, 
	portrait = :portrait
	WHERE character_id = :character_id`

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

	_, err = db.NamedExec(insertString, c)
	if err != nil {
		return err
	}

	return err
}

func UpdateCharacter(c CharacterRow) (err error) {
	db, err := GetDb()
	if err != nil {
		return err
	}

	_, err = db.NamedExec(updateString, c)

	return err
}
