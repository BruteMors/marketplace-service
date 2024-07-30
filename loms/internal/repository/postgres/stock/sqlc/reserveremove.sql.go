// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: reserveremove.sql

package sqlc

import (
	"context"
)

const reserveRemove = `-- name: ReserveRemove :exec
UPDATE items
SET
  reserved = GREATEST(reserved - data.count, 0),
  total_count = GREATEST(total_count - data.count, 0)
FROM (SELECT unnest($1::int[]) AS sku, unnest($2::int[]) AS count) AS data
WHERE items.sku = data.sku
`

type ReserveRemoveParams struct {
	Sku   []int32
	Count []int32
}

func (q *Queries) ReserveRemove(ctx context.Context, arg ReserveRemoveParams) error {
	_, err := q.db.Exec(ctx, reserveRemove, arg.Sku, arg.Count)
	return err
}
