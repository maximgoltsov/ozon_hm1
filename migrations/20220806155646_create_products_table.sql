-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.products (
    id          bigserial PRIMARY KEY,
    title       varchar(40) NOT NULL,
    price       integer NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.products;
-- +goose StatementEnd
