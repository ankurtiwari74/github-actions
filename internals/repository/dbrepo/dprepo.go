package dbrepo

import (
	"bookings/internals/config"
	"bookings/internals/repository"
	"database/sql"
)

type postgressDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgressRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgressDBRepo{
		App: a,
		DB:  conn,
	}
}
