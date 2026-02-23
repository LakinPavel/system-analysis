-- +goose Up
ALTER TABLE book ADD COLUMN booked BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE book DROP COLUMN booked;