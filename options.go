package postgres

import "time"

type Option func(*Postgres)

func ConnectionAttempts(count int) Option {
	return func(pg *Postgres) {
		pg.conectionAttempts = count
	}
}

func ConnectionTimeout(duration time.Duration) Option {
	return func(pg *Postgres) {
		pg.connectionTimeout = duration
	}
}

func MaxConnPoolSize(size int) Option {
	return func(pg *Postgres) {
		pg.maxConnPoolSize = size
	}
}
