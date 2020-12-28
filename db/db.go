package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gitlab.com/scalablespace/listener/app/models"
)

func NewDB(e models.Environment) (*sql.DB, error) {
	db, err := sql.Open("postgres", e.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(e.DBMaxOpenConns)
	db.SetMaxIdleConns(e.DBMaxIdleConns)
	db.SetConnMaxLifetime(e.DBConnMaxLifetime)
	return db, db.Ping()
}
