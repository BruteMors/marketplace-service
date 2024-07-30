package order

import (
	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/pg"
)

type Repository struct {
	db *pg.Client
}

func NewRepository(db *pg.Client) *Repository {
	repo := &Repository{
		db: db,
	}

	return repo
}
