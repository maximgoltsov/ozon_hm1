package db

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/maximgoltsov/botproject/internal/pkg/core/productType/models"
	repositoryPkg "github.com/maximgoltsov/botproject/internal/pkg/repository"
)

type productTypeRepository struct {
	pool *pgxpool.Pool
}

func NewProductTypeRepository(pool *pgxpool.Pool) repositoryPkg.ProductType {
	return &productTypeRepository{pool}
}

func (r *productTypeRepository) UpsertProductType(ctx context.Context, pt models.ProductType) (uint64, error) {
	if pt.Id != 0 {
		return pt.Id, updateProductType(ctx, &pt, r)
	} else {
		return pt.Id, addProductType(ctx, &pt, r)
	}
}

func addProductType(ctx context.Context, pt *models.ProductType, r *productTypeRepository) error {
	query, args, err := squirrel.Insert("product_types").
		Columns("name").
		Values(pt.Name).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(placeholder).
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

func updateProductType(ctx context.Context, pt *models.ProductType, r *productTypeRepository) error {
	query, args, err := squirrel.Update("product_types").
		Where(squirrel.Eq{
			"id": pt.Id,
		}).
		Set("name", pt.Name).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(placeholder).
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

func (r *productTypeRepository) DeleteProductTypeById(ctx context.Context, id uint64) error {
	query, args, err := squirrel.Delete("product_types").
		Where(squirrel.Eq{
			"id": id,
		}).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(placeholder).
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
func (r *productTypeRepository) GetProductType(ctx context.Context, id uint64) (models.ProductType, error) {
	query, args, err := squirrel.
		Select("id", "name").
		From("product_types").
		Where(squirrel.Eq{
			"id": id,
		}).
		PlaceholderFormat(placeholder).ToSql()

	if err != nil {
		return models.ProductType{}, err
	}

	var productType models.ProductType

	if err := pgxscan.Select(ctx, r.pool, &productType, query, args...); err != nil {
		return models.ProductType{}, err
	}

	return productType, nil
}
func (r *productTypeRepository) GetProductTypes(ctx context.Context, limit uint64, offset uint64, desc bool) []models.ProductType {
	querySelect := squirrel.Select("id", "name").From("product_types").PlaceholderFormat(placeholder)

	if limit != 0 {
		querySelect = querySelect.Offset(offset).Limit(limit)
	}

	if desc {
		querySelect = querySelect.OrderBy("id DESC")
	} else {
		querySelect = querySelect.OrderBy("id ASC")
	}

	query, args, err := querySelect.ToSql()

	var productTypes []models.ProductType

	if err != nil {
		return productTypes
	}

	if err := pgxscan.Select(ctx, r.pool, &productTypes, query, args...); err != nil {
		return productTypes
	}

	return productTypes
}
