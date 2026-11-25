// Package sqlite implementa core.Store con SQL directo sobre modernc.org/sqlite (Go puro, sin cgo)
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"runtime"
	"time"

	_ "modernc.org/sqlite"

	"captiveportal/core"
)

// maxReadConns acota el pool de lectura: más conexiones que núcleos no ayudan y cuestan memoria
const maxReadConns = 8

type Store struct {
	db  *sql.DB // escrituras: una sola conexión, SQLite tiene un solo escritor
	rdb *sql.DB // lecturas: varias conexiones; en WAL no esperan al escritor
}

func Open(path string) (*Store, error) {
	// synchronous=NORMAL: en WAL no hace fsync en cada commit (se pierde como mucho lo último
	// ante un corte de luz, no ante la caída del proceso). wal_autocheckpoint(4000) espacia el
	// checkpoint a ~16 MB para que no frene al escritor cada pocos segundos
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"+
		"&_pragma=synchronous(NORMAL)&_pragma=wal_autocheckpoint(4000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // el escritor nunca compite consigo mismo ni da SQLITE_BUSY
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	// query_only: el pool de lectura no puede escribir ni por accidente. Las migraciones ya corrieron
	rdb, err := sql.Open("sqlite", "file:"+path+"?_pragma=busy_timeout(5000)&_pragma=query_only(1)")
	if err != nil {
		db.Close()
		return nil, err
	}
	if n := min(runtime.NumCPU(), maxReadConns); n > 1 {
		rdb.SetMaxOpenConns(n)
		rdb.SetMaxIdleConns(n)
	}
	return &Store{db: db, rdb: rdb}, nil
}

func (s *Store) Close() error { return errors.Join(s.rdb.Close(), s.db.Close()) }

// Las horas se guardan en ISO 8601 UTC ('2026-11-19T18:03:00Z'), como exige el esquema
func iso(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func (s *Store) SaveVerified(ctx context.Context, v core.Visitor, at time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO verified (phone, name, id_number, company, verified_at) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (phone) DO UPDATE SET name = excluded.name, id_number = excluded.id_number,
			company = excluded.company, verified_at = excluded.verified_at`,
		v.Phone, v.Name, v.IDNumber, v.Company, iso(at))
	return err
}

func (s *Store) SaveAccess(ctx context.Context, a core.Access) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO access (ip, mac, phone, id_number, day, kind, start_at, end_at)
		VALUES (?, NULLIF(?, ''), ?, ?, ?, ?, ?, ?)`,
		a.IP, a.MAC, a.Phone, a.IDNumber, a.Day, a.Kind, iso(a.Start), iso(a.End))
	return err
}
