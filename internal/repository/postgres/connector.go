package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	panic("implement me")
}

func (d *Db) Get(ctx context.Context, key string) (json.RawMessage, error) {
	panic("implement me")
}

func (d *Db) Update(ctx context.Context, key string, metaData json.RawMessage) error {
	panic("implement me")
}

func (d *Db) Delete(ctx context.Context, key string) error {
	panic("implement me")
}