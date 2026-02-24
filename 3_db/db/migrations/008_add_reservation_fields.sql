-- +goose Up
ALTER TABLE book
    ADD COLUMN booked_by TEXT,
    ADD COLUMN reservation_start TIMESTAMP WITH TIME ZONE,
    ADD COLUMN reservation_end TIMESTAMP WITH TIME ZONE;

-- +goose Down
ALTER TABLE book
    DROP COLUMN booked_by,
    DROP COLUMN reservation_start,
    DROP COLUMN reservation_end;