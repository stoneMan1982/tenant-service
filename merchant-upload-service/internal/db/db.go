package db

import (
	"context"
	"database/sql"
	"fmt"
	"merchant-service/internal/config"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

const (
	DialectPostgres string = "postgres"
	DialectSQLite   string = "sqlite"
	DialectMySQL    string = "mysql"
)

// for sqlite, only support memory mode
func New(cfg *config.DBConfig) (*bun.DB, error) {
	var db *bun.DB
	switch cfg.Dialect {
	case DialectPostgres:
		db = newPostgresDB(cfg)
	case DialectSQLite:
		sqliteDB, err := newSQLiteDB()
		if err != nil {
			return nil, fmt.Errorf("failed to create SQLite connection: %w", err)
		}
		db = sqliteDB
	case DialectMySQL:
		mysqlDB, err := newMySQLDB(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create MySQL connection: %w", err)
		}
		db = mysqlDB
	default:
		return nil, fmt.Errorf("invalid dialect %s", cfg.Dialect)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func newPostgresDB(cfg *config.DBConfig) *bun.DB {
	// connection string: "postgres://username:password@host:port/database?sslmode=disable"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	configConnectionPool(sqlDB, cfg)
	return bun.NewDB(sqlDB, pgdialect.New())
}

// only support memory mode
func newSQLiteDB() (*bun.DB, error) {
	sqlDB, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(0)
	return bun.NewDB(sqlDB, sqlitedialect.New(), bun.WithDiscardUnknownColumns()), nil
}

func newMySQLDB(cfg *config.DBConfig) (*bun.DB, error) {
	// connection string: "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	configConnectionPool(sqlDB, cfg)
	return bun.NewDB(sqlDB, mysqldialect.New()), nil
}

func configConnectionPool(sqlDB *sql.DB, cfg *config.DBConfig) {
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime.Duration())
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime.Duration())
}

// isInMemoryDB checks if the DSN indicates an in-memory database
// TODO
func isInMemoryDB(dsn string) bool {
	return strings.Contains(dsn, ":memory:") || strings.Contains(dsn, "mode=memory")
}

func Close(db *bun.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
