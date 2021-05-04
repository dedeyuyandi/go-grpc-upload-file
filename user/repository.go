package user

import (
	"database/sql"

	"github.com/go-kit/kit/log"
)

type Repository interface {
}

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepo(db *sql.DB) Repository {
	return &repo{
		db: db,
	}
}
