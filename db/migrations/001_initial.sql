-- +goose Up
CREATE TABLE weapons (
    instance_id TEXT PRIMARY KEY,
    hash INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    power INTEGER NOT NULL,
    kills INTEGER NOT NULL DEFAULT 0,
    level INTEGER NOT NULL DEFAULT 0,
    location TEXT NOT NULL, -- 'Vault', 'Hunter', etc.
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE weapon_stats (
    instance_id TEXT NOT NULL,
    stat_name TEXT NOT NULL,
    value INTEGER NOT NULL,
    PRIMARY KEY (instance_id, stat_name),
    FOREIGN KEY (instance_id) REFERENCES weapons(instance_id) ON DELETE CASCADE
);

CREATE TABLE weapon_perks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    instance_id TEXT NOT NULL,
    perk_hash INTEGER NOT NULL,
    perk_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    FOREIGN KEY (instance_id) REFERENCES weapons(instance_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE weapon_perks;
DROP TABLE weapon_stats;
DROP TABLE weapons;
