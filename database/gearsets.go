package database

import (
	"encoding/json"

	"github.com/alyjay/xiv-be/types"
)

func InsertGearPiece(pieces []types.GearPieceRow, gsId string) (pieceRows []types.GearPieceRow, err error) {

	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	req, err := db.NamedQuery(`
			INSERT INTO gear_pieces (
				slot,
				source,
				have,
				augmented,
				priority
			) VALUES (
			 	:slot,
				:source,
				:have,
				:augmented,
				:priority
			 ) RETURNING *;
		`, pieces)
	if err != nil {
		return nil, err
	}
	for req.Next() {
		var id types.GearPieceRow
		req.StructScan(&id)
		_, err = db.Exec(`
			INSERT INTO gear_sets_gear_pieces (
				gear_set_id,
				gear_piece_id,
				gear_piece_slot
			) VALUES (
			 $1,
			 $2,
			 $3
			)
		`, gsId, id.Id, id.Slot)
		if err != nil {
			return nil, err
		}
		pieceRows = append(pieceRows, id)
	}

	return pieceRows, err
}

func InsertGearSetV2(g types.GearSetV2) (gs types.GearSetV2, err error) {

	db, err := GetDb()

	if err != nil {
		return g, err
	}
	gs = g

	_, err = db.NamedExec(`
		INSERT INTO gear_sets (
			id,
			user_id,
			character_id,
			name,
			job,
			index
		) VALUES (
		 :id,
		 :user_id,
		 :character_id,
		 :name,
		 :job,
		 :index
		 )
	`, g)

	if err != nil {
		return g, err
	}

	pieces := ConvertPieceMapToRow(g.Items)
	Items, err := InsertGearPiece(pieces, g.ID)

	for _, l := range Items {
		gs.Items[l.Slot] = l.GearPiece
	}

	return gs, err
}

func GetConfig(id string) (gs types.ItemsV2, err error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	query := db.QueryRow(`
		SELECT config from gear_sets WHERE id = $1
	`, id)

	var jsonData []byte
	err = query.Scan(&jsonData)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &gs)

	return gs, err
}

func ConvertPieceMapToRow(pieces types.ItemsV2) (rows []types.GearPieceRow) {
	for key, value := range pieces {
		rows = append(rows, types.GearPieceRow{
			GearPiece: value,
			Slot:      key,
		})

	}
	return rows
}

func SelectGearSetsForCharacterV2(id string, archived bool) (gs []types.GearSetV2, err error) {
	db, err := GetDb()

	if err != nil {
		return nil, err
	}

	err = db.Select(&gs, `
	SELECT 
		id, name, job
	FROM gear_sets WHERE
		character_id = $1
	AND
		archived = $2
	ORDER BY index ASC
	`, id, archived)

	if err != nil {
		return nil, err
	}

	var ret []types.GearSetV2

	for _, k := range gs {
		var Items []types.GearPieceRow
		err = db.Select(&Items, `
			SELECT 
				p.id, 
				p.source, 
				p.have, 
				p.augmented, 
				p.priority, 
				p.slot 
			FROM gear_pieces p
			JOIN gear_sets_gear_pieces r
			ON r.gear_piece_id = p.id
			WHERE r.gear_set_id = $1
		`, k.ID)
		k.Items = make(map[types.Slot]types.GearPiece)
		if len(Items) < 11 {
			k.Items, err = GetConfig(k.ID)
			pieces := ConvertPieceMapToRow(k.Items)
			InsertGearPiece(pieces, k.ID)
		} else {
			for _, l := range Items {
				k.Items[l.Slot] = l.GearPiece
			}
		}

		ret = append(ret, k)

	}

	return ret, err
}

func UpdateGearSetV2(g types.GearSetV2) (err error) {
	db, err := GetDb()

	if err != nil {
		return err
	}

	_, err = db.NamedExec(`
		UPDATE gear_sets SET
			name = :name,
			job = :job,
			index = :index,
			archived = :archived
		WHERE id = :id AND character_id = :character_id
	`, g)

	pieces := ConvertPieceMapToRow(g.Items)
	for _, piece := range pieces {
		_, err = db.NamedExec(`
			UPDATE gear_pieces SET
				have = :have,
				augmented = :augmented,
				priority = :priority
			WHERE id = :id AND slot = :slot
		`, piece)
		if err != nil {
			return err
		}
		_, err = db.Exec(`
			UPDATE gear_sets_gear_pieces SET
				gear_piece_id = $1
			WHERE 
				gear_set_id = $2 AND
				gear_piece_slot = $3
		`, piece.Id, g.ID, piece.Slot)
		if err != nil {
			return err
		}
	}

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
