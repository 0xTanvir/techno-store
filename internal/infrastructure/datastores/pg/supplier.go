package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"techno-store/internal/domain/bo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type supplierStore struct {
	dbPool *pgxpool.Pool
}

var supplierFields = []string{
	"id",
	"name",
	"email",
	"phone",
	"status_id",
	"is_verified_supplier",
	"created_at",
}

// supplierStore struct and its methods

// GetSupplierByID retrieves a supplier by its ID
func (s *supplierStore) GetSupplierByID(ctx context.Context, supplierID int64) (bo.Supplier, error) {
	var supplier bo.Supplier

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return supplier, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM suppliers WHERE id = $1", strings.Join(supplierFields, ","))
	row := conn.QueryRow(ctx, dbQuery, supplierID)

	if err = row.Scan(&supplier.ID, &supplier.Name, &supplier.Email, &supplier.Phone, &supplier.StatusID, &supplier.IsVerifiedSupplier, &supplier.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return bo.Supplier{}, bo.ErrSupplierNotFound
		}
		return bo.Supplier{}, err
	}

	return supplier, nil
}

// CreateSupplier inserts a new supplier into the database
func (s *supplierStore) CreateSupplier(ctx context.Context, supplier *bo.Supplier) error {
	insertMap := buildSupplierInsertMap(*supplier)
	if len(insertMap) < 1 {
		return fmt.Errorf("empty core insert for supplier")
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

	sqlQuery := fmt.Sprintf("INSERT INTO suppliers(%s)VALUES (%s) RETURNING id", strings.Join(fields, ","), strings.Join(placeholders, ","))

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var id sql.NullInt64
	if err := conn.QueryRow(ctx, sqlQuery, arguments...).Scan(&id); err != nil {
		return err
	}

	supplier.ID = id.Int64
	return nil
}

func buildSupplierInsertMap(i bo.Supplier) map[string]interface{} {
	insertedFields := make(map[string]interface{})

	for _, value := range supplierFields {
		switch value {
		case "name":
			if i.Name != "" {
				insertedFields[value] = i.Name
			}
		case "email":
			if i.Email != "" {
				insertedFields[value] = i.Email
			}
		case "phone":
			if i.Phone != "" {
				insertedFields[value] = i.Phone
			}
		case "status_id":
			insertedFields[value] = i.StatusID
		case "is_verified_supplier":
			insertedFields[value] = i.IsVerifiedSupplier
		}
	}

	return insertedFields
}

// UpdateSupplier updates a supplier in the database
func (s *supplierStore) UpdateSupplier(ctx context.Context, updateSupplier bo.SupplierUpdate) error {
	updateMap := buildSupplierUpdateMap(updateSupplier)
	if len(updateMap) < 1 {
		return errors.New("empty core update for supplier")
	}

	sqlQuery := "UPDATE suppliers SET "
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
	arguments = append(arguments, updateSupplier.ID)

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, sqlQuery, arguments...); err != nil {
		return err
	}

	return nil
}

func buildSupplierUpdateMap(u bo.SupplierUpdate) map[string]interface{} {
	updatedFields := make(map[string]interface{})

	if u.Name != nil {
		updatedFields["name"] = *u.Name
	}
	if u.Email != nil {
		updatedFields["email"] = *u.Email
	}
	if u.Phone != nil {
		updatedFields["phone"] = *u.Phone
	}
	if u.StatusID != nil {
		updatedFields["status_id"] = *u.StatusID
	}
	if u.IsVerifiedSupplier != nil {
		updatedFields["is_verified_supplier"] = *u.IsVerifiedSupplier
	}

	return updatedFields
}

// DeleteSupplier removes a supplier from the database
func (s *supplierStore) DeleteSupplier(ctx context.Context, supplierID int64) error {
	sqlQuery := `DELETE FROM suppliers WHERE id = $1`

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, sqlQuery, supplierID); err != nil {
		return err
	}

	return nil
}

// ListSuppliers retrieves a list of suppliers based on the query parameters
func (s *supplierStore) ListSuppliers(ctx context.Context, supplierQuery bo.SupplierQuery) (bo.PaginatedSupplierCollection, error) {
	pagingCollection := bo.PaginatedSupplierCollection{}

	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return pagingCollection, err
	}
	defer conn.Release()

	dbQuery := fmt.Sprintf("SELECT %s FROM suppliers ORDER BY name ASC LIMIT $1 OFFSET $2", strings.Join(supplierFields, ","))
	rows, err := conn.Query(ctx, dbQuery, supplierQuery.Limit, supplierQuery.Offset)
	if err != nil {
		return pagingCollection, err
	}
	defer rows.Close()

	var suppliers bo.SupplierCollection
	for rows.Next() {
		var supplier bo.Supplier
		if err := rows.Scan(&supplier.ID, &supplier.Name, &supplier.Email, &supplier.Phone, &supplier.StatusID, &supplier.IsVerifiedSupplier, &supplier.CreatedAt); err != nil {
			return pagingCollection, err
		}
		suppliers = append(suppliers, supplier)
	}

	if err = rows.Err(); err != nil {
		return pagingCollection, err
	}

	pagingCollection.Data = suppliers
	var totalRecord sql.NullInt64
	if err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM suppliers`).Scan(&totalRecord); err != nil {
		return pagingCollection, err
	}
	pagingCollection.Total = totalRecord.Int64

	return pagingCollection, nil
}
