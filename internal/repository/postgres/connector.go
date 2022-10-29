package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
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
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	var id int

	err = tx.QueryRowxContext(ctx, "INSERT INTO service (service, current_version) VALUES ($1, $2) RETURNING id", key, 1).Scan(&id)
	if err != nil {
		err = fmt.Errorf("incert service failed, error: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO config (service_id, metadata, version) VALUES ($1, $2, $3)", id, metaData, 1)
	if err != nil {
		err = fmt.Errorf("incert config failed, error: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit failed, error: %w", err)
	}

	return nil
}

func (d *Db) Get(ctx context.Context, key string) (json.RawMessage, error) {
	var answer json.RawMessage

	query := `
	SELECT c.metadata FROM config c
	JOIN service s ON c.service_id = s.id
	WHERE s.service = $1 AND c.version = s.current_version
	`

	err := d.db.GetContext(ctx, &answer, query, key)
	if err != nil {
		return nil, fmt.Errorf("get config failed, error: %w", err)
	}

	return answer, nil
}

func (d *Db) GetVersions(ctx context.Context, key string) ([]*domain.Config, error) {
	var result []*domain.Config

	query := `
	SELECT c.metadata as config, c.version as version FROM config c
	JOIN service s ON c.service_id = s.id
	WHERE s.service = $1
	`

	err := d.db.SelectContext(ctx, &result, query, key)
	if err != nil {
		return nil, fmt.Errorf("get versions failed, error: %w", err)
	}

	return result, nil
}

func (d *Db) Update(ctx context.Context, key string, metaData json.RawMessage) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO config (service, metadata) VALUES ($1, $2)", key, metaData)
	if err != nil {
		return err
	}
	/* Сделать методы апдейт и делит. Апдейт должен создавать новую версию в таблице конфиг и менять
	карэнт айди в таблице сервис. Делит должен удалять определенную версию конфига по переданному сервису и версии.
	Делит не может удалить карэнт версию.
	*/

	return nil
}

func (d *Db) Delete(ctx context.Context, key string, version int) error {
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
