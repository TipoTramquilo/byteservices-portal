package core

import (
	"context"
	"errors"
	"testing"
	"time"
)

// fakeStore guarda visitantes y accesos en memoria
type fakeStore struct {
	visitors []Visitor
	accesses []Access
}

func (f *fakeStore) SaveVerified(_ context.Context, v Visitor, _ time.Time) error {
	f.visitors = append(f.visitors, v)
	return nil
}

func (f *fakeStore) SaveAccess(_ context.Context, a Access) error {
	f.accesses = append(f.accesses, a)
	return nil
}

// fakeRouter registra cada concesión de internet
type fakeRouter struct {
	grants []time.Time
	ips    []string
}

func (f *fakeRouter) MAC(context.Context, string) (string, error) {
	return "52:54:00:77:30:65", nil
}

func (f *fakeRouter) GrantAccess(_ context.Context, ip string, until time.Time) error {
	f.ips = append(f.ips, ip)
	f.grants = append(f.grants, until)
	return nil
}

func newTestService(t *testing.T) (*Service, *fakeStore, *fakeRouter, time.Time) {
	t.Helper()
	loc, err := time.LoadLocation("America/Caracas")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 20, 10, 0, 0, 0, loc)
	store, router := &fakeStore{}, &fakeRouter{}
	svc := NewService(store, router, Config{Location: loc, ResetHour: 2, SessionTime: 30 * time.Minute})
	svc.now = func() time.Time { return now }
	return svc, store, router, now
}

var visitor = Visitor{
	Name:     "Ana Pérez",
	IDNumber: "V18456789",
	Company:  "Byte Services",
	Phone:    "+584125551234",
}

const visitorIP = "192.168.10.101"

func TestSendGrantsThirtyMinutes(t *testing.T) {
	svc, store, router, now := newTestService(t)

	if err := svc.Send(context.Background(), visitor, visitorIP); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if len(router.grants) != 1 {
		t.Fatalf("se esperaba 1 concesión, hubo %d", len(router.grants))
	}
	if want := now.Add(30 * time.Minute); !router.grants[0].Equal(want) {
		t.Fatalf("la sesión debe vencer a %v, venció %v", want, router.grants[0])
	}
	if router.ips[0] != visitorIP {
		t.Fatalf("concesión para %s, se esperaba %s", router.ips[0], visitorIP)
	}
	if len(store.visitors) != 1 || len(store.accesses) != 1 {
		t.Fatalf("se esperaba 1 visitante y 1 acceso, hubo %d y %d", len(store.visitors), len(store.accesses))
	}
}

func TestSendAgainGrantsNewSession(t *testing.T) {
	svc, _, router, _ := newTestService(t)
	for i := 0; i < 2; i++ {
		if err := svc.Send(context.Background(), visitor, visitorIP); err != nil {
			t.Fatalf("Send %d: %v", i, err)
		}
	}
	if len(router.grants) != 2 {
		t.Fatalf("cada envío debe dar una sesión nueva, hubo %d", len(router.grants))
	}
}

func TestSendRejectsInvalidData(t *testing.T) {
	svc, store, router, _ := newTestService(t)
	bad := visitor
	bad.Name = "Ana"
	if err := svc.Send(context.Background(), bad, visitorIP); !errors.Is(err, ErrInvalid) {
		t.Fatalf("se esperaba ErrInvalid, llegó %v", err)
	}
	if len(router.grants) != 0 || len(store.visitors) != 0 {
		t.Fatal("datos inválidos no deben dar internet ni guardar visitante")
	}
}
