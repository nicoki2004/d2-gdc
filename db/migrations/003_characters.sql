-- +goose Up
CREATE TABLE IF NOT EXISTS characters (
    character_id TEXT PRIMARY KEY,
    class_type INTEGER NOT NULL,
    light_level INTEGER NOT NULL,
    emblem_url TEXT,
    last_played DATETIME);

ALTER TABLE weapons ADD COLUMN tier TEXT;
ALTER TABLE weapons ADD COLUMN icon_url TEXT;
ALTER TABLE weapons ADD COLUMN slot TEXT;
ALTER TABLE weapons ADD COLUMN damage_type TEXT;
ALTER TABLE weapons ADD COLUMN ammo_type INTEGER;
ALTER TABLE weapons ADD COLUMN character_id TEXT;



-- +goose Down
DROP TABLE characters;
ALTER TABLE weapons DROP COLUMN tier TEXT;
ALTER TABLE weapons DROP COLUMN icon_url TEXT;
ALTER TABLE weapons DROP COLUMN slot TEXT;
ALTER TABLE weapons DROP COLUMN damage_type TEXT;
ALTER TABLE weapons DROP COLUMN ammo_type INTEGER;
ALTER TABLE weapons DROP COLUMN character_id TEXT;

