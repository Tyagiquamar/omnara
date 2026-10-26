-- +goose Up

ALTER TABLE machines
    ADD COLUMN delete_reason_code text,
    ADD COLUMN delete_reason_message text;
