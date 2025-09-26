package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgIdem struct {
	pool *pgxpool.Pool
}

func NewPostgresIdempotencyStore(ctx context.Context, dsn string) (IdempotencyStore, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	st := &pgIdem{pool: pool}
	if err := st.ensureTable(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return st, nil
}

func (s *pgIdem) ensureTable(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS idempotency_cache (
  key text PRIMARY KEY,
  req_hash text NOT NULL,
  status_code int NOT NULL,
  body bytea NOT NULL,
  expiry timestamptz NOT NULL
);
`)
	return err
}

func (s *pgIdem) Get(ctx context.Context, key string) (IdempotencyRecord, bool, error) {
	var rec IdempotencyRecord
	row := s.pool.QueryRow(ctx, `SELECT status_code, body, expiry, req_hash FROM idempotency_cache WHERE key=$1`, key)
	if err := row.Scan(&rec.StatusCode, &rec.Body, &rec.Expiry, &rec.ReqHash); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return IdempotencyRecord{}, false, err
		}
		return IdempotencyRecord{}, false, nil
	}
	return rec, true, nil
}

func (s *pgIdem) Set(ctx context.Context, key string, rec IdempotencyRecord) error {
	_, err := s.pool.Exec(ctx, `
INSERT INTO idempotency_cache(key, req_hash, status_code, body, expiry)
VALUES($1,$2,$3,$4,$5)
ON CONFLICT (key) DO UPDATE SET req_hash=EXCLUDED.req_hash, status_code=EXCLUDED.status_code, body=EXCLUDED.body, expiry=EXCLUDED.expiry
`, key, rec.ReqHash, rec.StatusCode, rec.Body, rec.Expiry)
	return err
}
