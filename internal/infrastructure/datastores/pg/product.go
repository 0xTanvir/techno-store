package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"techno-store/internal/domain/bo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productStore struct {
	dbPool *pgxpool.Pool
}

var productFields = []string{
	"id", "name", "description", "specifications", "brand_id",
	"category_id", "supplier_id", "unit_price", "discount_price",
	"tags", "status_id",
}

func (s *productStore) GetProductByID(ctx context.Context, productID int64) (bo.Product, error) {
	var (
		id             sql.NullInt64
		name           sql.NullString
		description    sql.NullString
		specifications sql.NullString
		brandID        sql.NullInt64
		categoryID     sql.NullInt64
		supplierID     sql.NullInt64
		unitPrice      sql.NullFloat64
		discountPrice  sql.NullFloat64
		tags           sql.NullString
		statusID       sql.NullInt64
	)

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return bo.Product{}, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM products WHERE id = $1", strings.Join(productFields, ","))
	row := conn.QueryRow(ctx, dbQuery, productID)

	if err = row.Scan(&id, &name, &description, &specifications, &brandID, &categoryID, &supplierID, &unitPrice, &discountPrice, &tags, &statusID); err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("product id does not exist", slog.Int64("id", productID))
			return bo.Product{}, bo.ErrProductNotFound
		}
		slog.Error("failed to scan product table row", "cause", err)
		return bo.Product{}, err
	}

	return bo.Product{
		ID:             id.Int64,
		Name:           name.String,
		Description:    description.String,
		Specifications: specifications.String,
		BrandID:        brandID.Int64,
		CategoryID:     categoryID.Int64,
		SupplierID:     supplierID.Int64,
		UnitPrice:      unitPrice.Float64,
		DiscountPrice:  discountPrice.Float64,
		Tags:           tags.String,
		StatusID:       statusID.Int64,
	}, nil
}

func (s *productStore) CreateProduct(ctx context.Context, product *bo.Product) error {
	insertMap := buildProductInsertMap(*product)
	if len(insertMap) < 1 {
		slog.Debug("empty core insert for product")
		return fmt.Errorf("empty core insert for product")
	}

	start := 1
	arguments := make([]interface{}, 0, len(insertMap))
	fields := []string{}
	placeholders := []string{}

	for field, v := range insertMap {
		fields = append(fields, field)
		arguments = append(arguments, v)
		placeholders = append(placeholders, "$"+strconv.Itoa(start))
		start++
	}

	sqlQuery := fmt.Sprintf("INSERT INTO products(%s) VALUES (%s) RETURNING id", strings.Join(fields, ","), strings.Join(placeholders, ","))

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var id sql.NullInt64
	if err := conn.QueryRow(ctx, sqlQuery, arguments...).Scan(&id); err != nil {
		return err
	}

	product.ID = id.Int64
	return nil
}

func buildProductInsertMap(p bo.Product) map[string]interface{} {
	insertedFields := make(map[string]interface{})

	// Assuming all fields except Description, Specifications, DiscountPrice, and Tags can be non-null.
	insertedFields["name"] = p.Name
	insertedFields["brand_id"] = p.BrandID
	insertedFields["category_id"] = p.CategoryID
	insertedFields["supplier_id"] = p.SupplierID
	insertedFields["unit_price"] = p.UnitPrice
	insertedFields["discount_price"] = p.DiscountPrice
	insertedFields["status_id"] = p.StatusID

	// Optional fields
	if p.Description != "" {
		insertedFields["description"] = p.Description
	}
	if p.Specifications != "" {
		insertedFields["specifications"] = p.Specifications
	}
	if p.Tags != "" {
		insertedFields["tags"] = p.Tags
	}

	return insertedFields
}

func (s *productStore) UpdateProduct(ctx context.Context, updateProduct bo.ProductUpdate) error {
	return WrapInTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		updateMap := buildProductUpdateMap(updateProduct)
		if len(updateMap) < 1 {
			slog.Debug("empty core update for product", slog.Int64("id", updateProduct.ID))
			return errors.New("empty core update for product")
		}

		sqlQuery := "UPDATE products SET "
		start := 1
		arguments := make([]interface{}, 0, len(updateMap)+1)

		for k, v := range updateMap {
			sqlQuery += k + "=$" + strconv.Itoa(start)
			arguments = append(arguments, v)
			if start < len(updateMap) {
				sqlQuery += ", "
			}
			start++
		}

		sqlQuery += fmt.Sprintf(" WHERE id = $%d", start)
		arguments = append(arguments, updateProduct.ID)

		commandTag, err := tx.Exec(ctx, sqlQuery, arguments...)
		if err != nil {
			slog.Error("failed to update product in database", "cause", err)
			return fmt.Errorf("failed to update product in database: %w", err)
		}

		if commandTag.RowsAffected() == 0 {
			slog.Warn("no rows affected when updating product", slog.Int64("productID", updateProduct.ID))
		}

		return nil
	})
}

func buildProductUpdateMap(u bo.ProductUpdate) map[string]interface{} {
	updateFields := map[string]interface{}{}

	if u.Name != nil {
		updateFields["name"] = *u.Name
	}
	if u.Description != nil {
		updateFields["description"] = *u.Description
	}
	if u.Specifications != nil {
		updateFields["specifications"] = *u.Specifications
	}
	if u.BrandID != nil {
		updateFields["brand_id"] = *u.BrandID
	}
	if u.CategoryID != nil {
		updateFields["category_id"] = *u.CategoryID
	}
	if u.SupplierID != nil {
		updateFields["supplier_id"] = *u.SupplierID
	}
	if u.UnitPrice != nil {
		updateFields["unit_price"] = *u.UnitPrice
	}
	if u.DiscountPrice != nil {
		updateFields["discount_price"] = *u.DiscountPrice
	}
	if u.Tags != nil {
		updateFields["tags"] = *u.Tags
	}
	if u.StatusID != nil {
		updateFields["status_id"] = *u.StatusID
	}

	return updateFields
}

func (s *productStore) DeleteProduct(ctx context.Context, productID int64) error {
	return WrapInTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		sqlProductQuery := `DELETE FROM products WHERE id = $1`
		if _, err := tx.Exec(ctx, sqlProductQuery, productID); err != nil {
			slog.Error("failed to delete product", slog.Int64("productID", productID), "cause", err)
			return err
		}

		return nil
	})
}

func (s *productStore) ListProducts(ctx context.Context, productQuery bo.ProductSearchQuery) (bo.PaginatedProductCollection, error) {
	pagingCollection := bo.PaginatedProductCollection{}

	dbQuery, countQuery := buildQuery(productQuery)

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return pagingCollection, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, dbQuery)
	if err != nil {
		slog.Error("failed to list products", "cause", err)
		return pagingCollection, err
	}
	defer rows.Close()

	var products bo.ProductCollection
	for rows.Next() {
		var (
			id             sql.NullInt64
			name           sql.NullString
			description    sql.NullString
			specifications sql.NullString
			brandID        sql.NullInt64
			categoryID     sql.NullInt64
			supplierID     sql.NullInt64
			unitPrice      sql.NullFloat64
			discountPrice  sql.NullFloat64
			tags           sql.NullString
			statusID       sql.NullInt64
		)

		if err := rows.Scan(&id, &name, &description, &specifications, &brandID, &categoryID, &supplierID, &unitPrice, &discountPrice, &tags, &statusID); err != nil {
			slog.Error("failed to scan product row", "cause", err)
			return pagingCollection, err
		}

		products = append(products, bo.Product{
			ID:             id.Int64,
			Name:           name.String,
			Description:    description.String,
			Specifications: specifications.String,
			BrandID:        brandID.Int64,
			CategoryID:     categoryID.Int64,
			SupplierID:     supplierID.Int64,
			UnitPrice:      unitPrice.Float64,
			DiscountPrice:  discountPrice.Float64,
			Tags:           tags.String,
			StatusID:       statusID.Int64,
		})
	}

	if err = rows.Err(); err != nil {
		slog.Error("failed during rows iteration", "cause", err)
		return pagingCollection, err
	}

	pagingCollection.Data = products
	var totalRecord sql.NullInt64
	if err = conn.QueryRow(ctx, countQuery).Scan(&totalRecord); err != nil {
		slog.Error("error scanning COUNT products row", "cause", err)
		return pagingCollection, err
	}
	pagingCollection.Total = totalRecord.Int64

	return pagingCollection, nil
}

func buildQuery(productQuery bo.ProductSearchQuery) (string, string) {
	var pFields = []string{
		"p.id", "p.name", "p.description", "p.specifications", "p.brand_id",
		"p.category_id", "p.supplier_id", "p.unit_price", "p.discount_price",
		"p.tags", "p.status_id",
	}

	sqlStatement := `SELECT %s FROM products p
		INNER JOIN brands b ON p.brand_id = b.id
		INNER JOIN categories c ON p.category_id = c.id
		INNER JOIN suppliers s ON p.supplier_id = s.id
		INNER JOIN product_stocks ps ON p.id = ps.product_id
		WHERE p.status_id = 1 AND ps.stock_quantity > 0`

	if productQuery.Filter.PriceRangeFilter.Min > 0 {
		sqlStatement += fmt.Sprintf(" AND p.unit_price >= %f", productQuery.Filter.PriceRangeFilter.Min)
	}
	if productQuery.Filter.PriceRangeFilter.Max > 0 {
		sqlStatement += fmt.Sprintf(" AND p.unit_price <= %f", productQuery.Filter.PriceRangeFilter.Max)
	}
	if len(productQuery.Filter.BrandFilter) > 0 {
		sqlStatement += fmt.Sprintf(" AND p.brand_id IN (%s)", strings.Trim(strings.Replace(fmt.Sprint(productQuery.Filter.BrandFilter), " ", ",", -1), "[]"))
	}
	if productQuery.Filter.CategoryFilter != 0 {
		sqlStatement += fmt.Sprintf(" AND p.category_id = %d", productQuery.Filter.CategoryFilter)
	}
	if productQuery.Filter.SupplierFilter != 0 {
		sqlStatement += fmt.Sprintf(" AND p.supplier_id = %d", productQuery.Filter.SupplierFilter)
	}
	if productQuery.Filter.VerifiedSupplierFilter {
		sqlStatement += " AND s.verified = true"
	}

	countQuery := fmt.Sprintf(sqlStatement, `COUNT(*)`)
	sqlStatement = fmt.Sprintf(sqlStatement, strings.Join(pFields, ","))
	if productQuery.Filter.Query != "" {
		sqlStatement += fmt.Sprintf(` AND p.name LIKE '%s'`, `%`+productQuery.Filter.Query+`%`)
		countQuery += fmt.Sprintf(` AND p.name LIKE '%s'`, `%`+productQuery.Filter.Query+`%`)
	}

	countQuery += ` LIMIT ` + strconv.Itoa(productQuery.Paging.Limit) + ` OFFSET ` + strconv.Itoa(productQuery.Paging.Offset)

	sqlStatement += ` ORDER BY ` + productQuery.Sort.Field + ` ` + productQuery.Sort.Order + ` LIMIT ` + strconv.Itoa(productQuery.Paging.Limit) + ` OFFSET ` + strconv.Itoa(productQuery.Paging.Offset)
	return sqlStatement, countQuery
}
