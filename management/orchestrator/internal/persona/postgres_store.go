package persona

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore persists Personas in PostgreSQL.
type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, databaseURL string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	s := &PostgresStore{pool: pool}
	if err := s.Migrate(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

func (s *PostgresStore) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *PostgresStore) Migrate(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS personas (
    id UUID PRIMARY KEY,
    version TEXT NOT NULL DEFAULT '0.1.0',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    display_name TEXT,
    data JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft',
    customer_owned BOOLEAN NOT NULL DEFAULT TRUE,
    template_id UUID,
    device_profile_id TEXT
);
CREATE INDEX IF NOT EXISTS idx_personas_status ON personas(status);
CREATE INDEX IF NOT EXISTS idx_personas_updated_at ON personas(updated_at DESC);
`)
	return err
}

func (s *PostgresStore) Create(ctx context.Context, p *Persona) (*Persona, error) {
	if p.ID == "" {
		p = New(p.DisplayName, p.Demographics.Location, p.Circadian.Timezone,
			p.Demographics.AgeRange[0], p.Demographics.AgeRange[1], p.Engagement.Type)
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.Version == "" {
		p.Version = "0.1.0"
	}
	if p.Status == "" {
		p.Status = StatusDraft
	}

	raw, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	_, err = s.pool.Exec(ctx, `
INSERT INTO personas (id, version, created_at, updated_at, display_name, data, status, customer_owned, template_id, device_profile_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`, p.ID, p.Version, p.CreatedAt, p.UpdatedAt, p.DisplayName, raw, string(p.Status), p.CustomerOwned, p.TemplateID, nullIfEmpty(p.DeviceProfileID))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PostgresStore) Get(ctx context.Context, id string) (*Persona, error) {
	var raw []byte
	err := s.pool.QueryRow(ctx, `SELECT data FROM personas WHERE id = $1`, id).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var p Persona
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PostgresStore) List(ctx context.Context) ([]*Persona, error) {
	rows, err := s.pool.Query(ctx, `SELECT data FROM personas ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Persona
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		var p Persona
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, err
		}
		result = append(result, &p)
	}
	return result, rows.Err()
}

func (s *PostgresStore) Update(ctx context.Context, p *Persona) error {
	p.UpdatedAt = time.Now().UTC()
	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `
UPDATE personas SET version=$2, updated_at=$3, display_name=$4, data=$5, status=$6,
    customer_owned=$7, template_id=$8, device_profile_id=$9
WHERE id=$1
`, p.ID, p.Version, p.UpdatedAt, p.DisplayName, raw, string(p.Status), p.CustomerOwned, p.TemplateID, nullIfEmpty(p.DeviceProfileID))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStore) Delete(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM personas WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

var ErrNotFound = errors.New("persona not found")
