package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"

	cachePkg "github.com/maximgoltsov/botproject/internal/pkg/core/product/cache"
	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) cachePkg.Interface {
	return &repository{pool}
}

func (r *repository) GetProducts(limit uint64, offset uint64, desc bool) []models.Product {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (r *repository) GetProduct(id uint64) (models.Product, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query, args, err := squirrel.
		Select("id", "title", "price", "type_id").
		From("products").
		Where(squirrel.Eq{
			"id": id,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return models.Product{}, err
	}

	var product models.Product

	if err := pgxscan.Select(ctx, r.pool, &product, query, args...); err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (r *repository) UpsertProduct(p models.Product) (uint64, error) {
	if p.Id != 0 {
		return p.Id, updateProduct(&p, r)
	} else {
		return p.Id, addProduct(&p, r)
	}
}

func addProduct(p *models.Product, r *repository) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query, args, err := squirrel.Insert("products").
		Columns("title", "price", "type_id").
		Values(p.Title, p.Price, p.Type_Id).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
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

func updateProduct(p *models.Product, r *repository) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query, args, err := squirrel.Update("products").
		Where(squirrel.Eq{
			"id": p.Id,
		}).
		Set("title", p.Title).
		Set("price", p.Price).
		Set("type_id", p.Type_Id).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
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

func (r *repository) DeleteProduct(p models.Product) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query, args, err := squirrel.Delete("products").
		Where(squirrel.Eq{
			"id": p.Id,
		}).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
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

func (r *repository) DeleteProductById(id uint64) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query, args, err := squirrel.Delete("products").
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
