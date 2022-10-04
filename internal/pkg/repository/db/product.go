package db

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
	repositoryPkg "github.com/maximgoltsov/botproject/internal/pkg/repository"
)

var placeholder = squirrel.Dollar

type productsRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) repositoryPkg.Product {
	return &productsRepository{pool}
}

func (r *productsRepository) GetProducts(ctx context.Context, limit uint64, offset uint64, desc bool) []models.Product {
	querySelect := squirrel.Select("id", "title", "price", "type_id").From("products")

	if limit != 0 {
		querySelect = querySelect.Offset(offset).Limit(limit)
	}

	if desc {
		querySelect = querySelect.OrderBy("id DESC")
	} else {
		querySelect = querySelect.OrderBy("id ASC")
	}

	query, args, err := querySelect.ToSql()

	var products []models.Product
	if err != nil {
		return products
	}

	if err := pgxscan.Select(ctx, r.pool, &products, query, args...); err != nil {
		return products
	}
	return products
}

func (r *productsRepository) GetProduct(ctx context.Context, id uint64) (models.Product, error) {
	query, args, err := squirrel.
		Select("id", "title", "price", "type_id").
		From("products").
		Where(squirrel.Eq{
			"id": id,
		}).PlaceholderFormat(placeholder).ToSql()

	if err != nil {
		return models.Product{}, err
	}

	var product models.Product

	if err := pgxscan.Select(ctx, r.pool, &product, query, args...); err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (r *productsRepository) UpsertProduct(ctx context.Context, p models.Product) (uint64, error) {
	if p.Id != 0 {
		return p.Id, updateProduct(ctx, &p, r)
	} else {
		return p.Id, addProduct(ctx, &p, r)
	}
}

func addProduct(ctx context.Context, p *models.Product, r *productsRepository) error {
	query, args, err := squirrel.Insert("products").
		Columns("title", "price", "type_id").
		Values(p.Title, p.Price, p.Type_Id).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(placeholder).
		ToSql()

	if err != nil {
		return err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&p.Id); err != nil {
		return err
	}

	return nil
}

func updateProduct(ctx context.Context, p *models.Product, r *productsRepository) error {
	query, args, err := squirrel.Update("products").
		Where(squirrel.Eq{
			"id": p.Id,
		}).
		Set("title", p.Title).
		Set("price", p.Price).
		Set("type_id", p.Type_Id).
		Suffix("RETURNING id").
		PlaceholderFormat(placeholder).
		ToSql()

	if err != nil {
		return err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return err
	}

	return nil
}

func (r *productsRepository) DeleteProduct(ctx context.Context, p models.Product) error {
	query, args, err := squirrel.Delete("products").
		Where(squirrel.Eq{
			"id": p.Id,
		}).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(placeholder).
		ToSql()

	if err != nil {
		return err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return err
	}

	return nil
}

func (r *productsRepository) DeleteProductById(ctx context.Context, id uint64) error {
	query, args, err := squirrel.Delete("products").
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
