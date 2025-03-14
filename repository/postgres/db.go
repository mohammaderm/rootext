package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type PostgresDB struct {
	conn   *sqlx.DB
	config Config
}

func (p PostgresDB) Conn() *sqlx.DB {
	return p.conn
}

func New(dbconfig Config) *PostgresDB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbconfig.Host, dbconfig.Port, dbconfig.Username, dbconfig.Password, dbconfig.Database)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("can not open postgres db: %w", err))
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &PostgresDB{
		conn:   db,
		config: dbconfig,
	}
}
