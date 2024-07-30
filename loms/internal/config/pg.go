package config

import (
	"errors"
	"os"
	"strings"
)

const (
	pgMasterDSN  = "PG_MASTER_DSN"
	pgReplicaDSN = "PG_REPLICA_DSN"
)

type PGConfig interface {
	MasterDSN() string
	ReplicaDSNs() []string
}

type pgConfig struct {
	masterDSN   string
	replicaDSNs []string
}

func NewPGConfig() (PGConfig, error) {
	masterDSN := os.Getenv(pgMasterDSN)
	if len(masterDSN) == 0 {
		return nil, errors.New("masterDSN not found")
	}

	replicaDSN := os.Getenv(pgReplicaDSN)
	var replicaDSNs []string
	if len(replicaDSN) > 0 {
		replicaDSNs = strings.Split(replicaDSN, ",")
	}

	return &pgConfig{
		masterDSN:   masterDSN,
		replicaDSNs: replicaDSNs,
	}, nil
}

func (cfg *pgConfig) MasterDSN() string {
	return cfg.masterDSN
}

func (cfg *pgConfig) ReplicaDSNs() []string {
	return cfg.replicaDSNs
}
