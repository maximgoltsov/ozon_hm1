-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.products
ADD COLUMN IF NOT EXISTS 
    type_id bigint REFERENCES public.product_types (id) NOT NULL
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.products
DROP COLUMN IF EXISTS 
    type_id
-- +goose StatementEnd