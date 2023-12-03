-- name: UpsertResource :one
INSERT INTO resources (
        id,
        name,
        project_id,
        description,
        tags,
        type,
        locality,
        json_data
    )
VALUES (
        :id,
        :name,
        :project_id,
        :description,
        :tags,
        :type,
        :locality,
        json(:data)
    ) ON CONFLICT (id, type) DO
UPDATE
SET name = excluded.name,
    project_id = excluded.project_id,
    description = excluded.description,
    tags = excluded.tags,
    type = excluded.type,
    locality = excluded.locality,
    json_data = excluded.json_data
WHERE resources.id = excluded.id
RETURNING resources.type,
    resources.json_data AS data;


-- name: DeleteResource :one
DELETE FROM resources
WHERE id = ?
RETURNING resources.id;