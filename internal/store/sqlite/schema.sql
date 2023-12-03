-- Made to be used with SQLite3
CREATE TABLE resources (
    id char(36) NOT NULL,
    name NOT NULL,
    project_id char(36) NOT NULL,
    description,
    tags,
    type int NOT NULL,
    locality char(100) NOT NULL,
    -- JSON data of the resource
    json_data NOT NULL,
    -- Some resources may share the same id, but have different types
    -- Example: a Cockpit in a project will use the same id as the project
    PRIMARY KEY (id, type)
);