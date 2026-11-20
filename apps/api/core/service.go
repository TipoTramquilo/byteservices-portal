package core

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// ErrInvalid es el único error de negocio: datos del formulario que no pasan la validación
var ErrInvalid = errors.New("invalid data")

type Config struct {
	Location    *time.Location // zona horaria del reset; solo para calcularlo, todo se guarda en UTC
	ResetHour   int            // hora local del reset diario (2 = 2:00 AM)
	ResetMinute int
	SessionTime time.Duration // duración de cada envío del formulario: 30 min
}

type Service struct {
	store  Store
	router Router
	cfg    Config
	now    func() time.Time
}

func NewService(store Store, router Router, cfg Config) *Service {
	return &Service{store: store, router: router, cfg: cfg, now: time.Now}
}

// Los mismos formatos que exige el esquema de la base
var (
	idNumberRe = regexp.MustCompile(`^[VE]\d{6,8}$`)
	phoneRe    = regexp.MustCompile(`^\+[1-9]\d{7,14}$`) // E.164
)

// Send recibe el formulario y abre internet al dispositivo por SessionTime. No hay código: cada
// envío da una sesión nueva, así que el visitante vuelve a llenar sus datos cuando vence
func (s *Service) Send(ctx context.Context, v Visitor, ip string) error {
	v.Name = strings.Join(strings.Fields(v.Name), " ")
	v.Company = strings.Join(strings.Fields(v.Company), " ")
	ip, ok := normalizeIP(ip)
	if len(strings.Fields(v.Name)) < 2 || utf8.RuneCountInString(v.Name) > 100 || !idNumberRe.MatchString(v.IDNumber) ||
		utf8.RuneCountInString(v.Company) < 2 || utf8.RuneCountInString(v.Company) > 100 ||
		!phoneRe.MatchString(v.Phone) || !ok {
		return ErrInvalid
	}
	now := s.now().In(s.cfg.Location)
	mac, err := s.mac(ctx, ip)
	if err != nil {
		return err
	}
	if err := s.store.SaveVerified(ctx, v, now); err != nil {
		return err
	}
	return s.grant(ctx, ip, mac, v, now, now.Add(s.cfg.SessionTime))
}

// mac le pide al router la MAC de esa IP y la deja como AA:BB:CC:DD:EE:FF; "" si no la conoce
func (s *Service) mac(ctx context.Context, ip string) (string, error) {
	mac, err := s.router.MAC(ctx, ip)
	if err != nil || mac == "" {
		return "", err
	}
	hw, err := net.ParseMAC(mac)
	if err != nil || len(hw) != 6 {
		return "", fmt.Errorf("router: MAC inválida %q para %s", mac, ip)
	}
	return strings.ToUpper(hw.String()), nil
}

func (s *Service) grant(ctx context.Context, ip, mac string, v Visitor, now, until time.Time) error {
	day, _ := s.window(now)
	if err := s.router.GrantAccess(ctx, ip, until); err != nil {
		return err
	}
	return s.store.SaveAccess(ctx, Access{IP: ip, MAC: mac, Phone: v.Phone, IDNumber: v.IDNumber, Day: day,
		Kind: Verified, Start: now, End: until})
}

// window devuelve el día que contiene now: va de un reset al siguiente. day es la fecha local
// en que empezó (a la 1:00 AM del 21 sigue siendo el día "2026-11-20") y end el próximo reset
func (s *Service) window(now time.Time) (day string, end time.Time) {
	y, m, d := now.Date()
	start := time.Date(y, m, d, s.cfg.ResetHour, s.cfg.ResetMinute, 0, 0, s.cfg.Location)
	if now.Before(start) {
		start = time.Date(y, m, d-1, s.cfg.ResetHour, s.cfg.ResetMinute, 0, 0, s.cfg.Location)
	}
	sy, sm, sd := start.Date()
	return start.Format(time.DateOnly), time.Date(sy, sm, sd+1, s.cfg.ResetHour, s.cfg.ResetMinute, 0, 0, s.cfg.Location)
}

// normalizeIP deja una sola forma por dirección: IPv4 mapeada en IPv6 pasa a IPv4
func normalizeIP(ip string) (string, bool) {
	a, err := netip.ParseAddr(ip)
	if err != nil {
		return "", false
	}
	return a.Unmap().WithZone("").String(), true
}
