package pg

import (
	"context"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Client struct {
	masterDBC   *pgxpool.Pool
	replicaDBCs []*pgxpool.Pool
}

func New(ctx context.Context, masterDSN string, replicaDSNs []string) (*Client, error) {
	masterDBC, err := pgxpool.New(ctx, masterDSN)
	if err != nil {
		return nil, errors.Errorf("failed to connect to master db: %v", err)
	}

	var replicaDBCs []*pgxpool.Pool
	for _, dsn := range replicaDSNs {
		dbc, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return nil, errors.Errorf("failed to connect to replica db: %v", err)
		}
		replicaDBCs = append(replicaDBCs, dbc)
	}

	return &Client{
		masterDBC:   masterDBC,
		replicaDBCs: replicaDBCs,
	}, nil
}

func (c *Client) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}
	for _, dbc := range c.replicaDBCs {
		dbc.Close()
	}
	return nil
}

func (c *Client) MasterDB() *pgxpool.Pool {
	return c.masterDBC
}

func (c *Client) ReplicaDB() *pgxpool.Pool {
	if len(c.replicaDBCs) == 0 {
		return c.masterDBC
	}
	selectedIndex := rand.Intn(len(c.replicaDBCs))
	return c.replicaDBCs[selectedIndex]
}

func (c *Client) ReplicaDBs() []*pgxpool.Pool {
	return c.replicaDBCs
}
