package core

import (
	"context"
	"time"
)

// Router le abre internet al dispositivo con esa IP hasta la hora indicada (portal cautivo
// de OPNsense, hotspot MikroTik, openNDS…). La IP es la de la petición, no la escribe el visitante
type Router interface {
	// MAC devuelve la MAC del dispositivo con esa IP (el router la sabe por ARP), o "" si no la conoce
	MAC(ctx context.Context, ip string) (string, error)
	GrantAccess(ctx context.Context, ip string, until time.Time) error
}
