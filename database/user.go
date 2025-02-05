package database

import (
	types "github.com/alyjay/xiv-be/types"
)

var insertCharacterString = `INSERT INTO users 
	(id, 
	username, 
	discord_id, 
	avatar, 
	accent_color, 
	auth_token, 
	expires) VALUES 
	(:id, 
	:username, 
	:discord_id, 
	:avatar, 
	:accent_color, 
	:auth_token, 
	:expires)
	ON CONFLICT (discord_id) DO UPDATE SET 
		accent_color = :accent_color, 
		avatar = :avatar, 
		auth_token = :auth_token, 
		expires = :expires
	`

var updateCharacterString = `UPDATE users SET 
		accent_color = :accent_color, 
		avatar = :avatar, 
		auth_token = :auth_token, 
		expires = :expires, 
		WHERE discord_id = :discord_id `

func InsertUser(u types.User) (id string, err error) {
	db, err := GetDb()
	if err != nil {
		return "", err
	}
	_, err = db.NamedExec(insertCharacterString, u)
	if err != nil {
		return "", err
	}
	err = db.Get(&u, `SELECT * FROM users WHERE discord_id=$1`, u.DiscordId)

	return u.ID, err
}

func UpdateUser(u types.User) (err error) {
	db, err := GetDb()
	if err != nil {
		return err
	}
	_, err = db.NamedExec(updateCharacterString, u)

	return err
}

func GetUser(id string) (user types.User, err error) {
	db, err := GetDb()
	if err != nil {
		return user, err
	}

	err = db.Get(&user, `SELECT * FROM users WHERE id=$1`, id)

	return user, err
}
