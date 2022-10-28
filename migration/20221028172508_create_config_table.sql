-- +goose Up
-- +goose StatementBegin
CREATE TABLE config
(
    service varchar NOT NULL UNIQUE ,
    metadata json NOT NULL,
    is_used boolean NOT NULL DEFAULT FALSE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE config
-- +goose StatementEnd
