package cart

import "sync"

type Repository struct {
	mutex sync.Mutex
	store map[int64]map[int64]uint16
}

func NewRepository() *Repository {
	return &Repository{
		store: make(map[int64]map[int64]uint16),
	}
}

type RepositoryWithMetrics struct {
	repo *Repository
}

func NewRepositoryWithMetrics(repo *Repository) *RepositoryWithMetrics {
	return &RepositoryWithMetrics{repo: repo}
}
