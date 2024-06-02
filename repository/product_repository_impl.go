package repository

import (
	"bubblevy/restful-api/helper"
	"bubblevy/restful-api/model/domain"
	"context"
	"database/sql"
	"errors"
)

type productRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &productRepositoryImpl{}
}

func (repository *productRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	query := "INSERT INTO products(product_name, price) VALUES (?, ?)"
	result, err := tx.ExecContext(ctx, query, product.ProductName, product.Price)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	product.Id = int(id)
	return product
}

func (repository *productRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	query := "UPDATE products SET product_name = ?, price = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, product.ProductName, product.Price, product.Id)
	helper.PanicIfError(err)

	return product
}

func (repository *productRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, product domain.Product) {
	query := "DELETE FROM products WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, product.Id)
	helper.PanicIfError(err)
}

func (repository *productRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, productId int) (domain.Product, error) {
	query := "SELECT id, product_name, price FROM products WHERE id = ?"
	rows, err := tx.QueryContext(ctx, query, productId)
	helper.PanicIfError(err)
	defer rows.Close()

	product := domain.Product{}
	if rows.Next() {
		err := rows.Scan(&product.Id, &product.ProductName, &product.Price)
		helper.PanicIfError(err)
		return product, nil
	} else {
		return product, errors.New("product not found")
	}
}

func (repository *productRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Product {
	query := "SELECT id, product_name, price FROM products"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		product := domain.Product{}
		err := rows.Scan(&product.Id, &product.ProductName, &product.Price)
		helper.PanicIfError(err)
		products = append(products, product)
	}
	return products
}
