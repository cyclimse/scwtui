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


-- name: ListTypedResourcesInProject :many
SELECT resources.type,
    resources.json_data AS data
FROM resources
WHERE project_id = ?
    AND type = ?
ORDER BY resources.name ASC;
