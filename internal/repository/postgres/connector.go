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
	d.db.MustBegin()
	d.db.MustExecContext(ctx, "INSERT INTO config (service, data) VALUES ($1, $2)", key, metaData)

	return nil
}

func (d *Db) Get(ctx context.Context, key string) (json.RawMessage, error) {
	var answer json.RawMessage
	d.db.MustBegin()
	err := d.db.GetContext(ctx, answer, "SELECT * FROM config WHERE service=$1", key)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (d *Db) Update(ctx context.Context, key string, metaData json.RawMessage) error {
	d.db.MustBegin()
	d.db.MustExecContext(ctx,"UPDATE config SET data = $1 WHERE service = $2", metaData, key)

	return nil
}

func (d *Db) Delete(ctx context.Context, key string) error {
	d.db.MustBegin()
	d.db.MustExecContext(ctx,"DELETE FROM config WHERE service = $1", key)

	return nil
}
