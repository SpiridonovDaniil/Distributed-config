package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SpiridonovDaniil/Distributed-config/internal/config"
	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Db struct {
	db *sqlx.DB
}

func New(cfg config.Postgres) *Db {
	conn, err := sqlx.Connect("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.User,
			cfg.Pass,
			cfg.Address,
			cfg.Port,
			cfg.Db,
		))
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

func (d *Db) RollBack(ctx context.Context, key string, version int) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	query := `
	SELECT s.service FROM config c
	JOIN service s ON c.service_id = s.id
	WHERE s.service = $1 AND c.version = $2 
	`

	var name string
	err = tx.GetContext(ctx, &name, query, key, version)
	if err != nil {
		err = fmt.Errorf("failed to get version config, error: %w", err)
		if !errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("the specified version does not exist")
		}

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE service SET current_version = $1 WHERE service = $2", version, key)
	if err != nil {
		err = fmt.Errorf("failed to change version config, error: %w", err)

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
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	var id int
	err = tx.GetContext(ctx, &id, "SELECT id FROM service WHERE service = $1", key)
	if err != nil {
		err = fmt.Errorf("failed select id config, err: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}

	query := `
	UPDATE service SET current_version = current_version + 1 
	WHERE service = $1 
	`

	_, err = tx.ExecContext(ctx, query, key)
	if err != nil {
		err = fmt.Errorf("failed change current version config, err: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}

	var version int
	err = tx.GetContext(ctx, &version, "SELECT current_version FROM service WHERE service = $1", key)
	if err != nil {
		err = fmt.Errorf("failed select new version config, err: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO config (service_id, metadata, version) VALUES ($1, $2, $3)", id, metaData, version)
	if err != nil {
		err = fmt.Errorf("failed create new version config, err: %w", err)

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

func (d *Db) Delete(ctx context.Context, key string, version int) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	var currentVersion int

	err = tx.GetContext(ctx, &currentVersion, "SELECT current_version FROM service WHERE service = $1", key)
	if err != nil {
		err = fmt.Errorf("get current version failed, error: %w", err)

		errTx := tx.Rollback()
		if errTx != nil {
			err = fmt.Errorf("rollback failed, error: %w", errTx)
		}

		return err
	}
	if currentVersion != version {
		_, err = tx.ExecContext(ctx, "DELETE FROM config WHERE version = $1", version)
		if err != nil {
			err = fmt.Errorf("delete config failed, error: %w", err)

			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("rollback failed, error: %w", errTx)
			}

			return err
		}
	} else {
		return fmt.Errorf("the distributed version cannot be removed")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
