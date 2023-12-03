-- name: GetResource :one
SELECT resources.type,
    resources.json_data AS data
FROM resources
WHERE id = ?;


-- name: ListAllResources :many
SELECT resources.type,
    resources.json_data AS data
FROM resources
ORDER BY resources.name ASC;
