-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    title VARCHAR(100),
    description VARCHAR(500),
    start_date BIGINT NOT NULL,
    end_date BIGINT NOT NULL,
    notify_before INT 
);

-- +goose Down
DROP TABLE events;