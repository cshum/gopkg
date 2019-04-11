package db

import (
	"context"
	"database/sql"
)

func Tx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	isSuccess := false
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if !isSuccess {
			_ = tx.Rollback()
		}
	}()
	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	isSuccess = true
	return nil
}
