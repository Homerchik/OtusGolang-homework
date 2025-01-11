-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY NOT NULL,
    userId UUID NOT NULL,
    title VARCHAR(100),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    description VARCHAR(500),
    notify_before INT 
);

-- +goose Down
DROP TABLE events;