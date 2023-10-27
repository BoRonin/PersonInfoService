package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil
}
func NewDB(dsn string) (*Postgres, error) {
	const op = "storage.postgres.NewDB"
	pg := &Postgres{}
	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	var counts int64
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgress not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			pg.Pool = connection
			return pg, nil
		}
		if counts > 10 {
			log.Println("Couldn't connect to Postgres in 20 seconds")
			return nil, fmt.Errorf("%s, %w", op, errors.New("couldn't connect to Postgres in 20 seconds"))
		}
		log.Println("Waiting for 2 seconds...")
		time.Sleep(time.Second * 2)
		continue
	}
}
