package postgres

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultConnectionAttempts = 5
	defaultConnectionTimeout  = time.Second
	defaultMaxConnPoolSize    = 1
)

// for a postgres singleton
var (
	postgresInstance *Postgres
	postgresOnce     sync.Once
	dbError          error
)

type Postgres struct {
	//how many times to try to connect to postgres before returning an error
	conectionAttempts int
	//how long to wait before trying to connect again
	connectionTimeout time.Duration
	//how many connections are allowed to exist at once
	maxConnPoolSize int
	//the connection pool
	ConnPool *pgxpool.Pool
	//sql statment builder
	Builder squirrel.StatementBuilderType
}

// NewOnce creates a postgres singleton
func NewOnce(dbUrl string, options ...Option) (*Postgres, error) {
	postgresOnce.Do(func() {
		postgresInstance = &Postgres{
			conectionAttempts: defaultConnectionAttempts,
			connectionTimeout: defaultConnectionTimeout,
			maxConnPoolSize:   defaultMaxConnPoolSize,
		}

		for _, option := range options {
			option(postgresInstance)
		}

		postgresInstance.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		poolConfig, err := pgxpool.ParseConfig(dbUrl)
		if err != nil {
			dbError = fmt.Errorf("postgres - NewOnce - pgxpool.ParseConfig: %w", err)
			return
		}
		poolConfig.MaxConns = int32(postgresInstance.maxConnPoolSize)

		for postgresInstance.conectionAttempts > 0 {
			postgresInstance.ConnPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
			if err == nil { //connected succesfully
				break
			}
			log.Printf("Trying to connect to Postgres, attempts left: %v\n", postgresInstance.conectionAttempts)
			time.Sleep(postgresInstance.connectionTimeout)
			postgresInstance.conectionAttempts--
		}

		if err != nil {
			dbError = fmt.Errorf("postgres - NewOnce - pgxpool.NewWithConfig: %w", err)
			return
		}
		slog.Info("succesfully connected to Postgres!")
	})
	if dbError != nil {
		return nil, dbError
	}
	return postgresInstance, nil
}

func New(dbUrl string, options ...Option) (*Postgres, error) {
	pgdb := &Postgres{
		conectionAttempts: defaultConnectionAttempts,
		connectionTimeout: defaultConnectionTimeout,
		maxConnPoolSize:   defaultMaxConnPoolSize,
	}

	for _, option := range options {
		option(pgdb)
	}

	pgdb.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}
	poolConfig.MaxConns = int32(pgdb.maxConnPoolSize)

	for pgdb.conectionAttempts > 0 {
		pgdb.ConnPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil { //connected succesfully
			break
		}
		log.Printf("Trying to connect to Postgres, attempts left: %v\n", pgdb.conectionAttempts)
		time.Sleep(pgdb.connectionTimeout)
		pgdb.conectionAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.NewWithConfig: %w", err)
	}
	slog.Info("succesfully connected to Postgres!")
	return pgdb, nil
}
