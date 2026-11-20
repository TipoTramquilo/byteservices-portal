package core

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Visitor struct {
	Name     string
	IDNumber string // cédula con su letra: V18456789
	Company  string // empresa donde trabaja
	Phone    string // E.164: +584125551234
}

// Code es el código pendiente de un teléfono. Solo se guarda su hash
type Code struct {
	Visitor
	Hash       string
	ExpiresAt  time.Time
	Attempts   int
	SentAt     time.Time
	Day        string // YYYY-MM-DD en la zona horaria del portal
	SendsToday int
}

// Store guarda visitantes, códigos pendientes y accesos a internet
type Store interface {
	IsVerified(ctx context.Context, phone string) (bool, error)
	SaveVerified(ctx context.Context, v Visitor, at time.Time) error
	GetCode(ctx context.Context, phone string) (Code, error) // ErrNotFound si no hay
	SaveCode(ctx context.Context, c Code) error              // reemplaza el anterior y reinicia los intentos
	AddAttempt(ctx context.Context, phone string) error
	DeleteCode(ctx context.Context, phone string) error
}
