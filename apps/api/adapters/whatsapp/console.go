// Package whatsapp implementa core.Sender
package whatsapp

import (
	"context"
	"log"
)

// Console imprime el código en vez de mandarlo, para desarrollar sin cuenta de WhatsApp
type Console struct{}

func (Console) SendCode(_ context.Context, phone, code string) error {
	log.Printf("[whatsapp] code for %s: %s", phone, code)
	return nil
}
