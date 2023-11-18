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

type productStockStore struct {
	dbPool *pgxpool.Pool
}

var productStockFields = []string{
	"id",
	"product_id",
	"stock_quantity",
	"updated_at",
}

func (s *productStockStore) GetProductStockByID(ctx context.Context, productStockID int64) (bo.ProductStock, error) {
	var (
		id            sql.NullInt64
		productID     sql.NullInt64
		stockQuantity sql.NullInt64
		updatedAt     sql.NullTime
	)

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return bo.ProductStock{}, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM product_stocks WHERE id = $1", strings.Join(productStockFields, ","))
	row := conn.QueryRow(ctx, dbQuery, productStockID)

	if err = row.Scan(&id, &productID, &stockQuantity, &updatedAt); err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("product stock id does not exist", slog.Int64("id", productStockID))
			return bo.ProductStock{}, bo.ErrProductStockNotFound
		}
		slog.Error("failed to scan product stock table row", "cause", err)
		return bo.ProductStock{}, err
	}

	return bo.ProductStock{
		ID:            id.Int64,
		ProductID:     productID.Int64,
		StockQuantity: stockQuantity.Int64,
		UpdatedAt:     updatedAt.Time,
	}, nil
}

func (s *productStockStore) CreateProductStock(ctx context.Context, productStock *bo.ProductStock) error {
	insertMap := buildProductStockInsertMap(*productStock)
	if len(insertMap) < 1 {
		slog.Debug("empty core insert for product stock")
		return fmt.Errorf("empty core insert for product stock")
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

	sqlQuery := fmt.Sprintf("INSERT INTO product_stocks(%s)VALUES (%s) RETURNING id", strings.Join(fields, ","), strings.Join(placeholders, ","))

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var id sql.NullInt64
	if err := conn.QueryRow(ctx, sqlQuery, arguments...).Scan(&id); err != nil {
		return err
	}

	productStock.ID = id.Int64
	return nil
}

func buildProductStockInsertMap(i bo.ProductStock) map[string]interface{} {
	insertedFields := make(map[string]interface{})

	for _, value := range productStockFields {
		switch value {
		case "product_id":
			insertedFields[value] = i.ProductID
		case "stock_quantity":
			insertedFields[value] = i.StockQuantity
		}
	}

	return insertedFields
}

func (s *productStockStore) UpdateProductStock(ctx context.Context, updateProductStock bo.ProductStockUpdate) error {
	return WrapInTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		updateMap := buildProductStockUpdateMap(updateProductStock)
		if len(updateMap) < 1 {
			slog.Debug("empty core update for product stock", slog.Int64("id", updateProductStock.ProductID))
			return errors.New("empty core update for product stock")
		}

		sqlQuery := "UPDATE product_stocks SET "
		start := 1
		arguments := make([]interface{}, 0, len(updateMap)+1)

		for k, v := range updateMap {
			sqlQuery = sqlQuery + k + "=$" + strconv.Itoa(start)
			arguments = append(arguments, v)
			if start < len(updateMap) {
				sqlQuery = sqlQuery + ", "
			}

			start++
		}

		sqlQuery = sqlQuery + fmt.Sprintf(" WHERE product_id = $%d", start)
		arguments = append(arguments, updateProductStock.ProductID)

		commandTag, err := tx.Exec(ctx, sqlQuery, arguments...)
		if err != nil {
			slog.Error("failed to update product stock in database", "cause", err)
			return fmt.Errorf("failed to update product stock in database: %w", err)
		}

		if commandTag.RowsAffected() == 0 {
			slog.Warn("no rows affected when updating product stock", slog.Int64("ProductID", updateProductStock.ProductID))
		}

		return nil
	})
}

func buildProductStockUpdateMap(u bo.ProductStockUpdate) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if u.StockQuantity != nil {
		updateFields["stock_quantity"] = *u.StockQuantity
	}

	return updateFields
}

func (s *productStockStore) DeleteProductStock(ctx context.Context, productStockID int64) error {
	return WrapInTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		sqlQuery := `DELETE FROM product_stocks WHERE id = $1`
		if commandTag, err := tx.Exec(ctx, sqlQuery, productStockID); err != nil {
			slog.Error("failed to delete product stock", slog.Int64("productStockID", productStockID), "cause", err)
			return err
		} else if commandTag.RowsAffected() == 0 {
			return bo.ErrProductStockNotFound
		}

		return nil
	})
}

func (s *productStockStore) ListProductStocks(ctx context.Context, productStockQuery bo.ProductStockQuery) (bo.PaginatedProductStockCollection, error) {
	pagingCollection := bo.PaginatedProductStockCollection{}

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return pagingCollection, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM product_stocks ORDER BY id ASC LIMIT $1 OFFSET $2", strings.Join(productStockFields, ","))
	rows, err := conn.Query(ctx, dbQuery, productStockQuery.Limit, productStockQuery.Offset)
	if err != nil {
		slog.Error("failed to list product stocks", "cause", err)
		return pagingCollection, err
	}
	defer rows.Close()

	var productStocks bo.ProductStockCollection
	for rows.Next() {
		var (
			id            sql.NullInt64
			productID     sql.NullInt64
			stockQuantity sql.NullInt64
			updatedAt     sql.NullTime
		)
		if err := rows.Scan(&id, &productID, &stockQuantity, &updatedAt); err != nil {
			slog.Error("failed to scan product stock row", "cause", err)
			return pagingCollection, err
		}
		productStocks = append(productStocks, bo.ProductStock{
			ID:            id.Int64,
			ProductID:     productID.Int64,
			StockQuantity: stockQuantity.Int64,
			UpdatedAt:     updatedAt.Time,
		})
	}

	if err = rows.Err(); err != nil {
		slog.Error("failed during rows iteration", "cause", err)
		return pagingCollection, err
	}

	pagingCollection.Data = productStocks
	var totalRecord sql.NullInt64
	if err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM product_stocks`).Scan(&totalRecord); err != nil {
		slog.Error("error scanning COUNT product stocks row", "cause", err)
		return pagingCollection, err
	}

	pagingCollection.Total = totalRecord.Int64
	return pagingCollection, nil
}
