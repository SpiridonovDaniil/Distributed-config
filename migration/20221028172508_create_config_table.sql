-- +goose Up
-- +goose StatementBegin
CREATE TABLE service
(
    id serial PRIMARY KEY,
    service varchar NOT NULL UNIQUE,
    current_version integer NOT NULL
);
CREATE TABLE config
(
    service_id integer REFERENCES service(id) NOT NULL,
    metadata json NOT NULL,
    version integer NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE config
-- +goose StatementEnd
