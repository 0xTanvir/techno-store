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

type categoryStore struct {
	dbPool *pgxpool.Pool
}

var categoryFields = []string{
	"id",
	"name",
	"parent_id",
	"sequence",
	"status_id",
	"created_at",
}

func (s *categoryStore) GetCategoryByID(ctx context.Context, categoryID int64) (bo.Category, error) {
	var (
		id        sql.NullInt64
		name      sql.NullString
		parentID  sql.NullInt64
		sequence  sql.NullInt64
		statusID  sql.NullInt64
		createdAt sql.NullTime
	)

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return bo.Category{}, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM categories WHERE id = $1", strings.Join(categoryFields, ","))
	row := conn.QueryRow(ctx, dbQuery, categoryID)

	if err = row.Scan(&id, &name, &parentID, &sequence, &statusID, &createdAt); err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("category id does not exist", slog.Int64("id", categoryID))
			return bo.Category{}, bo.ErrCategoryNotFound
		}
		slog.Error("failed to scan category table row", "cause", err)
		return bo.Category{}, err
	}

	return bo.Category{
		ID:        id.Int64,
		Name:      name.String,
		ParentID:  parentID.Int64,
		Sequence:  sequence.Int64,
		StatusID:  statusID.Int64,
		CreatedAt: createdAt.Time,
	}, nil
}

func (s *categoryStore) CreateCategory(ctx context.Context, category *bo.Category) error {
	insertMap := buildCategoryInsertMap(*category)
	if len(insertMap) < 1 {
		slog.Debug("empty insert for category")
		return fmt.Errorf("empty insert for category")
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

	sqlQuery := fmt.Sprintf("INSERT INTO categories(%s) VALUES (%s) RETURNING id", strings.Join(fields, ","), strings.Join(placeholders, ","))

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var id sql.NullInt64
	if err := conn.QueryRow(ctx, sqlQuery, arguments...).Scan(&id); err != nil {
		return err
	}

	category.ID = id.Int64
	return nil
}

func buildCategoryInsertMap(i bo.Category) map[string]interface{} {
	insertedFields := make(map[string]interface{})

	for _, value := range categoryFields {
		switch value {
		case "name":
			if i.Name != "" {
				insertedFields[value] = strings.ToLower(i.Name)
			}
		case "parent_id":
			insertedFields[value] = i.ParentID
		case "sequence":
			insertedFields[value] = i.Sequence
		case "status_id":
			insertedFields[value] = i.StatusID
		}
	}

	return insertedFields
}

func (s *categoryStore) UpdateCategory(ctx context.Context, updateCategory bo.CategoryUpdate) error {
	return WrapInTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		updateMap := buildCategoryUpdateMap(updateCategory)
		if len(updateMap) < 1 {
			slog.Debug("empty update for category", slog.Int64("id", updateCategory.ID))
			return errors.New("empty update for category")
		}

		sqlQuery := "UPDATE categories SET "
		start := 1
		arguments := make([]interface{}, 0, len(updateMap)+1)

		for k, v := range updateMap {
			sqlQuery += fmt.Sprintf("%s = $%d", k, start)
			arguments = append(arguments, v)
			if start < len(updateMap) {
				sqlQuery += ", "
			}
			start++
		}

		sqlQuery += fmt.Sprintf(" WHERE id = $%d", start)
		arguments = append(arguments, updateCategory.ID)

		commandTag, err := tx.Exec(ctx, sqlQuery, arguments...)
		if err != nil {
			slog.Error("failed to update category in database", "cause", err)
			return fmt.Errorf("failed to update category in database: %w", err)
		}

		if commandTag.RowsAffected() == 0 {
			slog.Warn("no rows affected when updating category", slog.Int64("categoryID", updateCategory.ID))
		}

		return nil
	})
}

func buildCategoryUpdateMap(u bo.CategoryUpdate) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if u.Name != nil {
		updateFields["name"] = *u.Name
	}
	if u.StatusID != nil {
		updateFields["status_id"] = *u.StatusID
	}
	if u.Sequence != nil {
		updateFields["sequence"] = *u.Sequence
	}

	return updateFields
}

func (s *categoryStore) DeleteCategory(ctx context.Context, categoryID int64) error {
	return WrapInTx(ctx, s.dbPool, func(tx pgx.Tx) error {
		sqlCategoryQuery := `DELETE FROM categories WHERE id = $1`
		if _, err := tx.Exec(ctx, sqlCategoryQuery, categoryID); err != nil {
			slog.Error("failed to delete category", slog.Int64("categoryID", categoryID), "cause", err)
			return err
		}
		return nil
	})
}

func (s *categoryStore) ListCategories(ctx context.Context) (bo.PaginatedCategoryCollection, error) {
	pagingCollection := bo.PaginatedCategoryCollection{}

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return pagingCollection, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM categories ORDER BY sequence ASC", strings.Join(categoryFields, ","))
	rows, err := conn.Query(ctx, dbQuery)
	if err != nil {
		slog.Error("failed to list categories", "cause", err)
		return pagingCollection, err
	}
	defer rows.Close()

	var categories bo.CategoryCollection
	for rows.Next() {
		var (
			id        sql.NullInt64
			name      sql.NullString
			parentID  sql.NullInt64
			sequence  sql.NullInt64
			statusID  sql.NullInt64
			createdAt sql.NullTime
		)
		if err := rows.Scan(&id, &name, &parentID, &sequence, &statusID, &createdAt); err != nil {
			slog.Error("failed to scan category row", "cause", err)
			return pagingCollection, err
		}

		categories = append(categories, bo.Category{
			ID:        id.Int64,
			Name:      name.String,
			ParentID:  parentID.Int64,
			Sequence:  sequence.Int64,
			StatusID:  statusID.Int64,
			CreatedAt: createdAt.Time,
		})
	}

	if err = rows.Err(); err != nil {
		slog.Error("failed during rows iteration", "cause", err)
		return pagingCollection, err
	}

	pagingCollection.Data = categories

	return pagingCollection, nil
}
