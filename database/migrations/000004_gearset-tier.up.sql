ALTER TABLE IF EXISTS gear_sets
    ADD COLUMN IF NOT EXISTS archived boolean default false,
    DROP CONSTRAINT gear_sets_character_id_fkey;

ALTER TABLE IF EXISTS gear_sets
    ADD CONSTRAINT gear_sets_character_id_fkey
        FOREIGN KEY (character_id)
        REFERENCES characters(character_id)
        ON DELETE CASCADE;

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS auth_token varchar(50) default '',
    ADD COLUMN IF NOT EXISTS expires bigint default 0;