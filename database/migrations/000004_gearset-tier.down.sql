ALTER TABLE IF EXISTS gear_sets
    DROP COLUMN IF EXISTS archived,
    DROP CONSTRAINT gear_sets_character_id_fkey;

ALTER TABLE IF EXISTS gear_sets
    ADD CONSTRAINT gear_sets_character_id_fkey
    FOREIGN KEY (character_id)
    REFERENCES characters(character_id);

ALTER TABLE IF EXISTS users
    DROP COLUMN IF EXISTS auth_token,
    DROP COLUMN IF EXISTS expires