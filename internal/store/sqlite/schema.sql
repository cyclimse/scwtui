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
    PRIMARY KEY (id)
);