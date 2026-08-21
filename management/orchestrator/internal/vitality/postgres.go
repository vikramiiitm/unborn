package vitality

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresTracker persists vitality scores.
type PostgresTracker struct {
	pool *pgxpool.Pool
}

func NewPostgresTracker(ctx context.Context, databaseURL string) (*PostgresTracker, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	t := &PostgresTracker{pool: pool}
	if err := t.migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return t, nil
}

func (t *PostgresTracker) migrate(ctx context.Context) error {
	_, err := t.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS vitality (
    persona_id UUID PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL DEFAULT 75,
    last_reason TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`)
	return err
}

func (t *PostgresTracker) Close() {
	if t.pool != nil {
		t.pool.Close()
	}
}

func (t *PostgresTracker) Get(ctx context.Context, personaID string) (*Score, error) {
	var s Score
	err := t.pool.QueryRow(ctx, `
SELECT persona_id::text, value, COALESCE(last_reason,''), updated_at FROM vitality WHERE persona_id = $1
`, personaID).Scan(&s.PersonaID, &s.Value, &s.LastReason, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return &Score{PersonaID: personaID, Value: DefaultScore, UpdatedAt: time.Now().UTC(), LastReason: "default"}, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (t *PostgresTracker) Ensure(ctx context.Context, personaID string) (*Score, error) {
	_, err := t.pool.Exec(ctx, `
INSERT INTO vitality (persona_id, value, last_reason, updated_at)
VALUES ($1, $2, 'initialized', NOW())
ON CONFLICT (persona_id) DO NOTHING
`, personaID, DefaultScore)
	if err != nil {
		return nil, err
	}
	return t.Get(ctx, personaID)
}

func (t *PostgresTracker) Adjust(ctx context.Context, personaID string, delta float64, reason string) (*Score, error) {
	if _, err := t.Ensure(ctx, personaID); err != nil {
		return nil, err
	}
	var s Score
	err := t.pool.QueryRow(ctx, `
UPDATE vitality SET
  value = GREATEST(0, LEAST(100, value + $2)),
  last_reason = $3,
  updated_at = NOW()
WHERE persona_id = $1
RETURNING persona_id::text, value, COALESCE(last_reason,''), updated_at
`, personaID, delta, reason).Scan(&s.PersonaID, &s.Value, &s.LastReason, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("adjust vitality: %w", err)
	}
	return &s, nil
}

func (t *PostgresTracker) List(ctx context.Context) ([]*Score, error) {
	rows, err := t.pool.Query(ctx, `
SELECT persona_id::text, value, COALESCE(last_reason,''), updated_at FROM vitality ORDER BY updated_at DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Score
	for rows.Next() {
		var s Score
		if err := rows.Scan(&s.PersonaID, &s.Value, &s.LastReason, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &s)
	}
	return out, rows.Err()
}
