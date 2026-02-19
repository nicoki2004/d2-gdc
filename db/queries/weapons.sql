-- name: UpsertWeapon :exec
INSERT INTO weapons (
    instance_id, hash, name, type, power, kills, level, location, 
    tier, icon_url, slot, damage_type, ammo_type, character_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(instance_id) DO UPDATE SET
    power = excluded.power,
    kills = excluded.kills,
    level = excluded.level,
    location = excluded.location,
    character_id = excluded.character_id;

-- name: ClearWeaponStats :exec
DELETE FROM weapon_stats WHERE instance_id = ?;
-- name: ClearWeaponPerks :exec
DELETE FROM weapon_perks WHERE instance_id = ?;

-- name: InsertWeaponStat :exec
INSERT INTO 
  weapon_stats (instance_id, stat_name, value) VALUES (?, ?, ?);

-- name: InsertWeaponPerk :exec
INSERT INTO weapon_perks (instance_id, perk_hash, perk_name, is_equipped, socket_index)
VALUES (?, ?, ?, ?, ?);

-- name: GetGodRollCandidates :many
SELECT 
  w.* FROM weapons w
JOIN weapon_perks p ON w.instance_id = p.instance_id
JOIN weapon_stats s ON w.instance_id = s.instance_id
WHERE p.perk_name = ? AND s.stat_name = 'Range' AND s.value > ?;

-- name: UpsertCharacter :exec
INSERT INTO characters (character_id, class_type, light_level, emblem_url, last_played)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT(character_id) DO UPDATE SET
    light_level = excluded.light_level,
    last_played = excluded.last_played;


-- name: ClearAllWeaponsData :exec
DELETE FROM weapons;

-- name: ClearAllWeaponStats :exec
DELETE FROM weapon_stats;

-- name: ClearAllWeaponsPerks :exec
DELETE FROM weapon_perks;
