-- +goose Up
ALTER TABLE weapon_perks ADD COLUMN is_equipped BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE weapon_perks ADD COLUMN socket_index INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE weapon_perks DROP COLUMN is_equipped;
ALTER TABLE weapon_perks DROP COLUMN socket_index;

