package transaction

import (
	"context"

	"github.com/BruteMors/marketplace-service/loms/pkg/client/db/pg"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type key string

const (
	TxKey key = "tx"
)

type TxManager interface {
	ReadCommitted(ctx context.Context, f func(ctx context.Context) error) error
}

type manager struct {
	db *pg.Client
}

func NewTransactionManager(db *pg.Client) TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) ReadCommitted(ctx context.Context, f func(context.Context) error) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn func(context.Context) error) (err error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.MasterDB().BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	ctx = m.makeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) makeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

type TxStarter interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func CreateTx(ctx context.Context, dbc TxStarter, options pgx.TxOptions) (
	tx pgx.Tx,
	commit func(ctx context.Context) error,
	rollback func(ctx context.Context) error,
	err error,
) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx, func(ctx context.Context) error { return nil }, func(ctx context.Context) error { return nil }, nil
	}

	beginTx, err := dbc.BeginTx(ctx, options)
	if err != nil {
		return nil, nil, nil, err
	}

	commit = func(ctx context.Context) error {
		return beginTx.Commit(ctx)
	}

	rollback = func(ctx context.Context) error {
		return beginTx.Rollback(ctx)
	}

	return beginTx, commit, rollback, nil
}

func CheckTx(ctx context.Context) (
	tx pgx.Tx, found bool,
) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx, true
	}

	return nil, false
}
