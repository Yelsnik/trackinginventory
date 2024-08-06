-- name: CreateInventory :one
INSERT INTO inventory (
  item,
  serial_number,
  price,
  owner
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetInventory :one
SELECT * FROM inventory
WHERE id = $1 LIMIT 1;

-- name: ListInventories :many
SELECT * FROM inventory
WHERE owner = $3
ORDER BY  id
LIMIT $1
OFFSET $2;

-- name: UpdateInventory :one
UPDATE inventory
  set item = $1,
   serial_number = $2
WHERE id = $3
RETURNING *;

-- name: DeleteInventory :exec
DELETE FROM inventory
WHERE id = $1;