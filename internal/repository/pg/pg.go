package pg

import (
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(dbClient *bun.DB) *Repository {
	return &Repository{
		db: dbClient,
	}
}
