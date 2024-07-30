package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/BruteMors/marketplace-service/loms/internal/config"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

var (
	pgMasterDSN   = ""
	pgReplicaDSNs = []string{}
	migrationsDir = "./migrations"
)

var (
	conn *pgxpool.Pool
)

func setupTest(t *testing.T) {
	err := config.Load("../.env")
	require.NoError(t, err)

	pgConfig, err := config.NewPGConfig()
	require.NoError(t, err)

	pgMasterDSN = pgConfig.MasterDSN()
	pgReplicaDSNs = pgConfig.ReplicaDSNs()

	db, err := sql.Open("pgx", pgMasterDSN)
	require.NoError(t, err)

	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	err = goose.Up(db, migrationsDir)
	require.NoError(t, err)

	conn, err = pgxpool.New(context.Background(), pgMasterDSN)
	require.NoError(t, err)
}

func teardownTest(t *testing.T) {
	db, err := sql.Open("pgx", pgMasterDSN)
	require.NoError(t, err)

	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	err = goose.Down(db, migrationsDir)
	require.NoError(t, err)

	if conn != nil {
		conn.Close()
	}
}
