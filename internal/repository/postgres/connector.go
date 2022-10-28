package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Db struct {
	db *sqlx.DB
}

func New(
	address string,
	user string,
	pass string,
	db string,
) *Db {
	conn, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, address, db))
	if err != nil {
		log.Fatal(err)
	}

	return &Db{db: conn}
}

func (d *Db) Create(ctx context.Context, key string, metaData json.RawMessage) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO config (service, metadata) VALUES ($1, $2)", key, metaData)
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) Get(ctx context.Context, key string) (json.RawMessage, error) {
	var answer json.RawMessage
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}

	err = sqlx.GetContext(ctx, tx, &answer, "SELECT metadata FROM config WHERE service = $1", key)
	if err != nil {
		err = fmt.Errorf("get metadata failed, error: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return nil, err
	}

	_, err = tx.ExecContext(ctx, "UPDATE config SET is_used = $1 WHERE service = $2", true, key)
	if err != nil {
		err = fmt.Errorf("update is_used failed, error: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return nil, err
	}

	return answer, nil
}

func (d *Db) Update(ctx context.Context, key string, metaData json.RawMessage) error {
	_, err := d.db.ExecContext(ctx, "UPDATE config SET metadata = $1 WHERE service = $2", metaData, key)
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) Delete(ctx context.Context, key string) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	var isUsed bool

	err = sqlx.GetContext(ctx, tx, &isUsed, "SELECT is_used FROM config WHERE service = $1", key)
	if err != nil {
		err = fmt.Errorf("get isUsed failed, error: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM config WHERE service = $1", key)
	if err != nil {
		err = fmt.Errorf("delete config failed, error: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
