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

type AccessKind string

const Verified AccessKind = "verified" // sesión de 30 min abierta al enviar el formulario

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

// Store guarda visitantes y los accesos a internet
type Store interface {
	SaveVerified(ctx context.Context, v Visitor, at time.Time) error
	SaveAccess(ctx context.Context, a Access) error
}
