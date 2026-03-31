package app

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBBundle struct {
	Writer *sqlx.DB
	Reader *sqlx.DB
}

func OpenDatabases(cfg Config) (*DBBundle, error) {
	writer, err := sqlx.Open("mysql", cfg.WriterDSN)
	if err != nil {
		return nil, err
	}
	configurePool(writer, cfg)

	readerDSN := cfg.ReaderDSN
	if readerDSN == "" {
		readerDSN = cfg.WriterDSN
	}

	reader, err := sqlx.Open("mysql", readerDSN)
	if err != nil {
		_ = writer.Close()
		return nil, err
	}
	configurePool(reader, cfg)

	return &DBBundle{Writer: writer, Reader: reader}, nil
}

func configurePool(db *sqlx.DB, cfg Config) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(5 * time.Minute)
}
