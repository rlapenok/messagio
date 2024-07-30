package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rlapenok/messagio/internal/utils"
)

func createTableMessages(db *sqlx.DB) {
	var err error
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id UUID PRIMARY KEY, message TEXT,processed BOOLEAN DEFAULT false)"); err != nil {
		utils.Logger.Fatal(err.Error())
	}
}

func withTransaction[T any](ctx context.Context, r repository, query func(context context.Context, tx *sqlx.Tx) (T, error)) (T, error) {
	var t T
	//create tx opts
	opts := sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelDefault,
	}
	//start tx
	tx, err := r.db.BeginTxx(ctx, &opts)
	if err != nil {
		return t, err
	}
	//handle err when return from func
	defer func(tx *sqlx.Tx) {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				errors.Join(err, rollbackErr)
			}
		}
	}(tx)
	//execute query
	result, err := query(ctx, tx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return result, err
	}
	if err := tx.Commit(); err != nil {
		utils.Logger.Error(err.Error())
		return t, err
	}
	return result, nil

}
