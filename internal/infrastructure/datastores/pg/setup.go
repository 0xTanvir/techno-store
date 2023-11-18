package pg

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"techno-store/config"
	"techno-store/internal/domain/definition"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	rwInstanceOnce sync.Once
	rwInstance     *pgxpool.Pool
)

// products, product_stocks.
func GetInstance(config *config.DBConfig) definition.DataStore {
	dbpool := RwInstance(config)

	return definition.DataStore{
		Brand:        &brandStore{dbPool: dbpool},
		Category:     &categoryStore{dbPool: dbpool},
		Supplier:     &supplierStore{dbPool: dbpool},
		Product:      &productStore{dbPool: dbpool},
		ProductStock: &productStockStore{dbPool: dbpool},
	}
}

func RwInstance(config *config.DBConfig) *pgxpool.Pool {
	rwInstanceOnce.Do(func() {
		rwInstance = getInstance(config)
	})
	return rwInstance
}

func Close(dbpool *pgxpool.Pool) {
	dbpool.Close()
}

func getInstance(config *config.DBConfig) *pgxpool.Pool {
	var (
		dbpool *pgxpool.Pool
		err    error
	)

	url := fmt.Sprintf("postgres://%s:%s@%s/%s", config.User, config.Password, config.Address, config.Database)

	dbpool, err = pgxpool.New(context.Background(), url)
	if err != nil {
		slog.Error("Unable to connect to database", err)
	}

	return dbpool
}

// WrapInTx starts a transaction on the given database connection pool, calls the
// given transaction function with a pointer to the transaction, and either commits
// or rolls back the transaction based on whether an error occurred or not. If the
// transaction function panics, the transaction is rolled back and the panic is
// propagated. The function returns an error if the transaction could not be started
// or committed/rolled back.
func WrapInTx(ctx context.Context, dbpool *pgxpool.Pool, txFunc func(pgx.Tx) error) error {
	tx, err := dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				slog.Error("problem during panic roll back", "cause", err)
			}
			panic(p)
		} else if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				slog.Error("problem during error roll back", "cause", err)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = txFunc(tx)
	return err
}
