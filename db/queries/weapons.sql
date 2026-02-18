-- name: UpsertWeapon :exec
INSERT INTO 
  weapons (instance_id, hash, name, type, power, kills, level, location, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(instance_id) DO UPDATE SET
    power = excluded.power,
    kills = excluded.kills,
    level = excluded.level,
    location = excluded.location,
    updated_at = CURRENT_TIMESTAMP;

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
