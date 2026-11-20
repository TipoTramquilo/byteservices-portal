package core

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStore struct {
	verified map[string]Visitor
	codes    map[string]Code
}

func newFakeStore() *fakeStore {
	return &fakeStore{verified: map[string]Visitor{}, codes: map[string]Code{}}
}

func (f *fakeStore) IsVerified(_ context.Context, phone string) (bool, error) {
	_, ok := f.verified[phone]
	return ok, nil
}

func (f *fakeStore) SaveVerified(_ context.Context, v Visitor, _ time.Time) error {
	f.verified[v.Phone] = v
	return nil
}

func (f *fakeStore) GetCode(_ context.Context, phone string) (Code, error) {
	c, ok := f.codes[phone]
	if !ok {
		return Code{}, ErrNotFound
	}
	return c, nil
}

func (f *fakeStore) SaveCode(_ context.Context, c Code) error {
	f.codes[c.Phone] = c
	return nil
}

func (f *fakeStore) AddAttempt(_ context.Context, phone string) error {
	c := f.codes[phone]
	c.Attempts++
	f.codes[phone] = c
	return nil
}

func (f *fakeStore) DeleteCode(_ context.Context, phone string) error {
	delete(f.codes, phone)
	return nil
}

type fakeSender struct {
	calls int
	code  string
}

func (f *fakeSender) SendCode(_ context.Context, _, code string) error {
	f.calls++
	f.code = code
	return nil
}

func newTestService(t *testing.T) (*Service, *fakeStore, *fakeSender, *time.Time) {
	t.Helper()
	loc, err := time.LoadLocation("America/Caracas")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 20, 10, 0, 0, 0, loc)
	store, sender := newFakeStore(), &fakeSender{}
	svc := NewService(store, sender, Config{
		Location: loc, ResetHour: 2, CodeTTL: 10 * time.Minute, ResendWait: time.Minute,
		MaxSendsPerDay: 5, MaxAttempts: 5, Secret: []byte("test-secret"),
	})
	svc.now = func() time.Time { return now }
	return svc, store, sender, &now
}

var visitor = Visitor{Name: "Ana Pérez", IDNumber: "V18456789", Company: "Byte Services", Phone: "+584125551234"}

func TestSendManda_UnCodigo(t *testing.T) {
	svc, _, sender, _ := newTestService(t)
	verified, err := svc.Send(context.Background(), visitor)
	if err != nil || verified {
		t.Fatalf("Send: verified=%v err=%v", verified, err)
	}
	if sender.calls != 1 || len(sender.code) != 6 {
		t.Fatalf("se esperaba 1 código de 6 dígitos, llamadas=%d código=%q", sender.calls, sender.code)
	}
}

func TestSendMuyPronto_NoManda(t *testing.T) {
	svc, _, sender, _ := newTestService(t)
	ctx := context.Background()
	if _, err := svc.Send(ctx, visitor); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Send(ctx, visitor); !errors.Is(err, ErrTooSoon) {
		t.Fatalf("se esperaba ErrTooSoon, llegó %v", err)
	}
	if sender.calls != 1 {
		t.Fatalf("no debía mandar otro código, llamadas=%d", sender.calls)
	}
}

func TestVerifyCorrecto_VerificaYBorraElCodigo(t *testing.T) {
	svc, store, sender, _ := newTestService(t)
	ctx := context.Background()
	if _, err := svc.Send(ctx, visitor); err != nil {
		t.Fatal(err)
	}
	if err := svc.Verify(ctx, visitor.Phone, sender.code); err != nil {
		t.Fatalf("Verify con código correcto: %v", err)
	}
	if _, ok := store.verified[visitor.Phone]; !ok {
		t.Fatal("el teléfono debe quedar verificado")
	}
	if _, ok := store.codes[visitor.Phone]; ok {
		t.Fatal("el código usado debe borrarse")
	}
}

func TestVerifyIncorrecto_NoVerifica(t *testing.T) {
	svc, store, sender, _ := newTestService(t)
	ctx := context.Background()
	if _, err := svc.Send(ctx, visitor); err != nil {
		t.Fatal(err)
	}
	wrong := "000000"
	if wrong == sender.code {
		wrong = "111111"
	}
	if err := svc.Verify(ctx, visitor.Phone, wrong); !errors.Is(err, ErrWrongCode) {
		t.Fatalf("se esperaba ErrWrongCode, llegó %v", err)
	}
	if _, ok := store.verified[visitor.Phone]; ok {
		t.Fatal("un código incorrecto no debe verificar")
	}
}

func TestSendYaVerificado_NoMandaCodigo(t *testing.T) {
	svc, store, sender, now := newTestService(t)
	store.verified[visitor.Phone] = visitor
	_ = now
	verified, err := svc.Send(context.Background(), visitor)
	if err != nil || !verified {
		t.Fatalf("un verificado debe responder verified=true: verified=%v err=%v", verified, err)
	}
	if sender.calls != 0 {
		t.Fatal("un verificado no debe recibir código")
	}
}

func TestSendDatosInvalidos(t *testing.T) {
	svc, _, sender, _ := newTestService(t)
	bad := visitor
	bad.Name = "Ana"
	if _, err := svc.Send(context.Background(), bad); !errors.Is(err, ErrInvalid) {
		t.Fatalf("se esperaba ErrInvalid, llegó %v", err)
	}
	if sender.calls != 0 {
		t.Fatal("datos inválidos no deben mandar código")
	}
}
