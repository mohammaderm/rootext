package migrator

import (
	"database/sql"
	"fmt"

	"github.com/mohammaderm/rootext/repository/postgres"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dbConfig   postgres.Config
	dialect    string
	migrations *migrate.FileMigrationSource
}

func New(dbConfig postgres.Config) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/postgres/migrations",
	}
	return Migrator{
		dbConfig:   dbConfig,
		dialect:    "postgres",
		migrations: migrations,
	}
}

func (m Migrator) Up() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Database))
	if err != nil {
		panic(fmt.Errorf("can not open postgre db: %w", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("can not apply migration: %w", err))
	}
	fmt.Printf("applied %d migrations! \n", n)
}

func (m Migrator) Down() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Database))
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("can't rollback migrations: %v", err))
	}
	fmt.Printf("Rollbacked %d migrations!\n", n)
}
