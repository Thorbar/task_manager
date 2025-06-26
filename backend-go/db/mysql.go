package db

import (
	"database/sql"
	"fmt"
	"task-manager/backend-go/config"
	"task-manager/backend-go/internal/i18n"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// ConnectDB initializes the MySQL connection using environment variables
func ConnectDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("db.open"), err)

	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("%s: %w", i18n.T("db.ping"), err)
	}

	return nil
}
