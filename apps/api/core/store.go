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

type AccessKind string

const (
	Free     AccessKind = "free"     // la hora gratis para que llegue el código
	Verified AccessKind = "verified" // hasta el próximo reset, tras verificar
)

type Access struct {
	IP       string
	MAC      string // vacía si el router no la conoce
	Phone    string
	IDNumber string
	Day      string
	Kind     AccessKind
	Start    time.Time
	End      time.Time
}

// Store guarda visitantes, códigos pendientes y accesos a internet
type Store interface {
	IsVerified(ctx context.Context, phone string) (bool, error)
	SaveVerified(ctx context.Context, v Visitor, at time.Time) error
	GetCode(ctx context.Context, phone string) (Code, error) // ErrNotFound si no hay
	SaveCode(ctx context.Context, c Code) error              // reemplaza el anterior y reinicia los intentos
	AddAttempt(ctx context.Context, phone string) error
	DeleteCode(ctx context.Context, phone string) error
	HadFreeAccess(ctx context.Context, ip, mac, phone, idNumber, day string) (bool, error) // por cualquiera de ellos; mac vacía no cuenta
	SaveAccess(ctx context.Context, a Access) error
}
