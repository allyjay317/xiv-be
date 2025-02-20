CREATE TABLE
    IF NOT EXISTS gear_pieces (
        id serial primary key,
        slot INTEGER,
        source INTEGER,
        have BOOLEAN,
        augmented BOOLEAN,
        priority INTEGER
    );

CREATE TABLE
    IF NOT EXISTS gear_sets_gear_pieces (
        gear_set_id uuid REFERENCES gear_sets (id) ON DELETE CASCADE,
        gear_piece_id integer REFERENCES gear_pieces (id) ON DELETE CASCADE
    );