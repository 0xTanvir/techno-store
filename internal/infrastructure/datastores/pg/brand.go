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

type brandStore struct {
	dbPool *pgxpool.Pool
}

var brandFields = []string{
	"id",
	"name",
	"status_id",
	"created_at",
}

func (s *brandStore) GetBrandByID(ctx context.Context, brandID int64) (bo.Brand, error) {
	var (
		id        sql.NullInt64
		name      sql.NullString
		statusID  sql.NullInt64
		createdAt sql.NullTime
	)

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return bo.Brand{}, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM brands WHERE id = $1", strings.Join(brandFields, ","))
	row := conn.QueryRow(ctx, dbQuery, brandID)

	if err = row.Scan(&id, &name, &statusID, &createdAt); err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("brand id does not exist", slog.Int64("id", brandID))
			return bo.Brand{}, bo.ErrBrandNotFound
		}
		slog.Error("failed to scan brand table row", "cause", err)
		return bo.Brand{}, err
	}

	return bo.Brand{
		ID:        id.Int64,
		Name:      name.String,
		StatusID:  statusID.Int64,
		CreatedAt: createdAt.Time,
	}, nil
}

func (s *brandStore) CreateBrand(ctx context.Context, brand *bo.Brand) error {
	insertMap := buildBrandInsertMap(*brand)
	if len(insertMap) < 1 {
		slog.Debug("empty core insert for brand")
		return fmt.Errorf("empty core insert for brand")
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

	sqlQuery := fmt.Sprintf("INSERT INTO brands(%s)VALUES (%s) RETURNING id", strings.Join(fields, ","), strings.Join(placeholders, ","))

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var id sql.NullInt64
	if err := conn.QueryRow(ctx, sqlQuery, arguments...).Scan(&id); err != nil {
		return err
	}

	brand.ID = id.Int64
	return nil
}

func buildBrandInsertMap(i bo.Brand) map[string]interface{} {
	insertedFields := make(map[string]interface{})

	for _, value := range brandFields {
		switch value {
		case "name":
			if i.Name != "" {
				insertedFields[value] = strings.ToLower(i.Name)
			}
		case "status_id":
			insertedFields[value] = i.StatusID
		}
	}

	return insertedFields
}

func (s *brandStore) UpdateBrand(ctx context.Context, updateBrand bo.BrandUpdate) error {
	return WrapInTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		updateMap := buildBrandUpdateMap(updateBrand)
		if len(updateMap) < 1 {
			slog.Debug("empty core update for brand", slog.Int64("id", updateBrand.ID))
			return errors.New("empty core update for brand")
		}

		sqlQuery := "UPDATE brands SET "
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

		sqlQuery = sqlQuery + fmt.Sprintf(" WHERE id = $%d", start)
		arguments = append(arguments, updateBrand.ID)

		commandTag, err := tx.Exec(ctx, sqlQuery, arguments...)
		if err != nil {
			slog.Error("failed to update brand in database", "cause", err)
			return fmt.Errorf("failed to update brand in database: %w", err)
		}

		if commandTag.RowsAffected() == 0 {
			slog.Warn("no rows affected when update brand", slog.Int64("brandID", updateBrand.ID))
		}

		return nil
	})
}

func buildBrandUpdateMap(u bo.BrandUpdate) map[string]interface{} {
	updateFields := []string{
		"name",
		"status_id",
	}
	updatedFields := make(map[string]interface{})

	for _, value := range updateFields {
		switch value {
		case "name":
			if u.Name != nil {
				updatedFields[value] = *u.Name
			}
		case "status_id":
			if u.StatusID != nil {
				updatedFields[value] = *u.StatusID
			}
		}
	}

	return updatedFields
}

func (s *brandStore) DeleteBrand(ctx context.Context, brandID int64) error {
	return WrapInTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		sqlBrandQuery := `DELETE FROM brands WHERE id = $1`
		if _, err := tx.Exec(ctx, sqlBrandQuery, brandID); err != nil {
			slog.Error("failed to delete brand", slog.Int64("brandID", brandID), "cause", err)
			return err
		}

		return nil
	})
}

func (s *brandStore) ListBrands(ctx context.Context, brandQuery bo.BrandQuery) (bo.PaginatedBrandCollection, error) {
	pagingCollection := bo.PaginatedBrandCollection{}

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return pagingCollection, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM brands ORDER BY name ASC LIMIT $1 OFFSET $2", strings.Join(brandFields, ","))
	rows, err := conn.Query(ctx, dbQuery, brandQuery.Limit, brandQuery.Offset)
	if err != nil {
		slog.Error("failed to list brands", "cause", err)
		return pagingCollection, err
	}
	defer rows.Close()

	var brands bo.BrandCollection
	for rows.Next() {
		var (
			id        sql.NullInt64
			name      sql.NullString
			statusID  sql.NullInt64
			createdAt sql.NullTime
		)
		if err := rows.Scan(&id, &name, &statusID, &createdAt); err != nil {
			slog.Error("failed to scan brand row", "cause", err)
			return pagingCollection, err
		}
		brands = append(brands, bo.Brand{
			ID:        id.Int64,
			Name:      name.String,
			StatusID:  statusID.Int64,
			CreatedAt: createdAt.Time,
		})
	}

	if err = rows.Err(); err != nil {
		slog.Error("failed during rows iteration", "cause", err)
		return pagingCollection, err
	}

	pagingCollection.Data = brands
	var (
		totalRecord sql.NullInt64
	)
	if err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM brands`).Scan(&totalRecord); err != nil {
		slog.Error("error scanning COUNT brands row", "cause", err)
		return pagingCollection, err
	}

	pagingCollection.Total = totalRecord.Int64

	return pagingCollection, nil
}
