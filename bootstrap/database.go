package bootstrap

import (
	"database/sql"
	"log"

	"wallet-api/config"

	_ "github.com/go-sql-driver/mysql"
)

// NewDatabase creates a new database connection using the provided config
func NewDatabase(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, logError("failed to open database", err)
	}

	// Apply connection pool settings from config
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, logError("failed to ping database", err)
	}

	// Log pool stats for monitoring
	log.Printf("Database connection pool initialized: MaxOpen=%d, MaxIdle=%d, ConnMaxLifetime=%v, ConnMaxIdleTime=%v",
		cfg.MaxOpenConns,
		cfg.MaxIdleConns,
		cfg.ConnMaxLifetime,
		cfg.ConnMaxIdleTime)

	return db, nil
}

func logError(message string, err error) error {
	log.Println(message + ": " + err.Error())
	return err
}
