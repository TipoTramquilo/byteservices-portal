package core

import "context"

// Sender le hace llegar el código al visitante (WhatsApp)
type Sender interface {
	SendCode(ctx context.Context, phone, code string) error
}
