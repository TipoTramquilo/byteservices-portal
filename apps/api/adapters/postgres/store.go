// Package postgres implementa core.Store sobre PostgreSQL (Supabase u otro Postgres), con pgx
// como driver de database/sql
package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"captiveportal/core"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Store struct {
	db *sql.DB
}

// Open conecta con la URL de Postgres (DATABASE_URL) y aplica las migraciones pendientes
func Open(ctx context.Context, url string) (*Store, error) {
	// Modo de consulta simple: el pooler de transacciones de Supabase reutiliza conexiones, y las
	// sentencias preparadas de pgx se chocan entre ellas (SQLSTATE 42P05). Se fuerza aquí, para que
	// no dependa de la URL
	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_URL: %w", err)
	}
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	db := stdlib.OpenDB(*cfg)
	db.SetMaxOpenConns(10)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("postgres: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

// migrate aplica en orden cada archivo de migrations/ que todavía no está registrado
func (s *Store) migrate(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		return err
	}
	names, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)
	for _, name := range names {
		version, err := strconv.Atoi(strings.SplitN(strings.TrimPrefix(name, "migrations/"), "_", 2)[0])
		if err != nil {
			return fmt.Errorf("migración %s: nombre inválido", name)
		}
		var done int
		if err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations WHERE version = $1`, version).Scan(&done); err != nil {
			return err
		}
		if done > 0 {
			continue
		}
		body, err := migrations.ReadFile(name)
		if err != nil {
			return err
		}
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, string(body)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migración %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) IsVerified(ctx context.Context, phone string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM verified WHERE phone = $1`, phone).Scan(&n)
	return n > 0, err
}

func (s *Store) SaveVerified(ctx context.Context, v core.Visitor, at time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO verified (phone, name, id_number, company, verified_at) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (phone) DO UPDATE SET name = EXCLUDED.name, id_number = EXCLUDED.id_number,
			company = EXCLUDED.company, verified_at = EXCLUDED.verified_at`,
		v.Phone, v.Name, v.IDNumber, v.Company, at.UTC())
	return err
}

func (s *Store) GetCode(ctx context.Context, phone string) (core.Code, error) {
	var c core.Code
	err := s.db.QueryRowContext(ctx, `
		SELECT phone, name, id_number, company, hash, expires_at, attempts, sent_at, day, sends_today
		FROM codes WHERE phone = $1`, phone).
		Scan(&c.Phone, &c.Name, &c.IDNumber, &c.Company, &c.Hash, &c.ExpiresAt, &c.Attempts, &c.SentAt, &c.Day, &c.SendsToday)
	if errors.Is(err, sql.ErrNoRows) {
		return c, core.ErrNotFound
	}
	return c, err
}

func (s *Store) SaveCode(ctx context.Context, c core.Code) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO codes (phone, name, id_number, company, hash, expires_at, attempts, sent_at, day, sends_today)
		VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8, $9)
		ON CONFLICT (phone) DO UPDATE SET name = EXCLUDED.name, id_number = EXCLUDED.id_number,
			company = EXCLUDED.company, hash = EXCLUDED.hash, expires_at = EXCLUDED.expires_at, attempts = 0,
			sent_at = EXCLUDED.sent_at, day = EXCLUDED.day, sends_today = EXCLUDED.sends_today`,
		c.Phone, c.Name, c.IDNumber, c.Company, c.Hash, c.ExpiresAt.UTC(), c.SentAt.UTC(), c.Day, c.SendsToday)
	return err
}

func (s *Store) AddAttempt(ctx context.Context, phone string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE codes SET attempts = attempts + 1 WHERE phone = $1`, phone)
	return err
}

func (s *Store) DeleteCode(ctx context.Context, phone string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM codes WHERE phone = $1`, phone)
	return err
}
