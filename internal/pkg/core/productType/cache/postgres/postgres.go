package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"

	cachePkg "github.com/maximgoltsov/botproject/internal/pkg/core/productType/cache"
	"github.com/maximgoltsov/botproject/internal/pkg/core/productType/models"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) cachePkg.Interface {
	return &repository{pool}
}

func (r *repository) UpsertProductType(ctx context.Context, pt models.ProductType) (uint64, error) {
	if pt.Id != 0 {
		return pt.Id, updateProductType(ctx, &pt, r)
	} else {
		return pt.Id, addProductType(ctx, &pt, r)
	}
}

func addProductType(ctx context.Context, pt *models.ProductType, r *repository) error {
	query, args, err := squirrel.Insert("product_types").
		Columns("name").
		Values(pt.Name).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&pt.Id); err != nil {
		return err
	}

	return nil
}

func updateProductType(ctx context.Context, pt *models.ProductType, r *repository) error {
	query, args, err := squirrel.Update("product_types").
		Where(squirrel.Eq{
			"id": pt.Id,
		}).
		Set("name", pt.Name).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(); err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteProductTypeById(ctx context.Context, id uint64) error {
	query, args, err := squirrel.Delete("product_types").
		Where(squirrel.Eq{
			"id": id,
		}).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(); err != nil {
		return err
	}

	return nil
}
func (r *repository) GetProductType(ctx context.Context, id uint64) (models.ProductType, error) {
	query, args, err := squirrel.
		Select("id", "name").
		From("product_types").
		Where(squirrel.Eq{
			"id": id,
		}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return models.ProductType{}, err
	}

	var productType models.ProductType

	if err := pgxscan.Select(ctx, r.pool, &productType, query, args...); err != nil {
		return models.ProductType{}, err
	}

	return productType, nil
}
func (r *repository) GetProductTypes(ctx context.Context) []models.ProductType {
	query, args, err := squirrel.
		Select("id", "name").
		From("product_types").
		OrderBy("id ASC").
		ToSql()

	var productTypes []models.ProductType

	if err != nil {
		return productTypes
	}

	if err := pgxscan.Select(ctx, r.pool, &productTypes, query, args...); err != nil {
		return productTypes
	}

	return productTypes
}
