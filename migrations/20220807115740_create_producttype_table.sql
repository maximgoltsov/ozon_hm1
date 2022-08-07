-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.product_types (
    id         bigserial PRIMARY KEY,
    name       varchar(40) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.product_types;
-- +goose StatementEnd
