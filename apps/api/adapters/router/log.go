// Package router implementa core.Router. Solo existe la versión de log hasta que se
// confirme el equipo (mini PC con OPNsense o Raspberry)
package router

import (
	"context"
	"log"
	"time"
)

// Log registra el acceso sin abrir nada
type Log struct{}

// MAC no la conoce: sin router, el acceso queda solo con IP, teléfono y cédula
func (Log) MAC(context.Context, string) (string, error) { return "", nil }

func (Log) GrantAccess(_ context.Context, ip string, until time.Time) error {
	log.Printf("[router] grant %s until %s", ip, until.UTC().Format(time.RFC3339))
	return nil
}
