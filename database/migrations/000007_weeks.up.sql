CREATE TABLE
    IF NOT EXISTS weeks (
        week_start date,
        character_id varchar(50) REFERENCES characters (character_id),
        tomes integer default 0,
        floor_one boolean default FALSE,
        floor_two boolean default FALSE,
        floor_three boolean default FALSE,
        floor_four boolean default FALSE,
        PRIMARY KEY (character_id, week_start)
    );

CREATE TABLE
    IF NOT EXISTS purchases (
        id serial primary key,
        character_id varchar(50),
        week_start date,
        slot int,
        purchase_type VARCHAR(50) not null,
        FOREIGN KEY (character_id, week_start) REFERENCES weeks (character_id, week_start)
    );