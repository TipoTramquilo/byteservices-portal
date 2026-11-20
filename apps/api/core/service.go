package core

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"net"
	"net/netip"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// Errores que el adaptador HTTP traduce a códigos de estado
var (
	ErrInvalid   = errors.New("invalid data")
	ErrTooSoon   = errors.New("code sent too recently")
	ErrTooMany   = errors.New("too many sends or attempts")
	ErrWrongCode = errors.New("wrong or expired code")
)

type Config struct {
	Location       *time.Location // zona horaria del reset; solo para calcularlo, todo se guarda en UTC
	ResetHour      int            // hora local del reset diario (2 = 2:00 AM)
	ResetMinute    int
	FreeTime       time.Duration // internet al enviar, para que llegue el código por WhatsApp
	CodeTTL        time.Duration
	ResendWait     time.Duration // igual a WAIT en el portal
	MaxSendsPerDay int
	MaxAttempts    int
	Secret         []byte // clave HMAC de los hashes de código
}

type Service struct {
	store  Store
	sender Sender
	router Router
	cfg    Config
	now    func() time.Time
	locks  [phoneLocks]sync.Mutex
}

// phoneLocks son las franjas del lock por teléfono: dos números pueden compartir franja, lo que
// solo cuesta algo de espera entre ellos. Con 256 basta para que el choque sea raro
const phoneLocks = 256

// lockPhone bloquea ese teléfono y devuelve la función que lo suelta
func (s *Service) lockPhone(phone string) func() {
	h := fnv.New32a()
	h.Write([]byte(phone))
	m := &s.locks[h.Sum32()%phoneLocks]
	m.Lock()
	return m.Unlock
}

func NewService(store Store, sender Sender, router Router, cfg Config) *Service {
	return &Service{store: store, sender: sender, router: router, cfg: cfg, now: time.Now}
}

// Los mismos formatos que exige el esquema de la base
var (
	idNumberRe = regexp.MustCompile(`^[VE]\d{6,8}$`)
	phoneRe    = regexp.MustCompile(`^\+[1-9]\d{7,14}$`) // E.164
	codeRe     = regexp.MustCompile(`^\d{6}$`)
)

// Send inicia una visita. Un teléfono verificado otro día recibe internet hasta el
// próximo reset y verified=true; si no, se manda un código por WhatsApp y el dispositivo recibe
// su hora gratis, una vez al día por IP, MAC, teléfono y cédula: los teléfonos rotan la MAC
// (y con ella la IP) y el teléfono y la cédula los escribe el visitante, ninguno basta solo
func (s *Service) Send(ctx context.Context, v Visitor, ip string) (verified bool, err error) {
	v.Name = strings.Join(strings.Fields(v.Name), " ")
	v.Company = strings.Join(strings.Fields(v.Company), " ")
	ip, ok := normalizeIP(ip)
	if len(strings.Fields(v.Name)) < 2 || utf8.RuneCountInString(v.Name) > 100 || !idNumberRe.MatchString(v.IDNumber) ||
		utf8.RuneCountInString(v.Company) < 2 || utf8.RuneCountInString(v.Company) > 100 ||
		!phoneRe.MatchString(v.Phone) || !ok {
		return false, ErrInvalid
	}
	now := s.now().In(s.cfg.Location)

	// el chequeo y el guardado del código van juntos: si dos peticiones del mismo teléfono se
	// cruzan, ambas ven el código anterior y cada una manda el suyo. El lock se suelta antes de
	// llamar a Twilio, que es lo lento y no toca el estado
	unlock := s.lockPhone(v.Phone)
	verified, code, err := s.prepare(ctx, v, ip, now)
	unlock()
	if err != nil || verified {
		return verified, err
	}
	return false, s.sender.SendCode(ctx, v.Phone, code)
}

// prepare decide si toca enviar y guarda el código; devuelve el código a mandar (vacío si no hay
// que mandar nada). Se llama con el lock del teléfono tomado
func (s *Service) prepare(ctx context.Context, v Visitor, ip string, now time.Time) (verified bool, code string, err error) {
	mac, err := s.mac(ctx, ip)
	if err != nil {
		return false, "", err
	}
	verified, err = s.store.IsVerified(ctx, v.Phone)
	if err != nil || verified {
		if verified {
			_, end := s.window(now)
			err = s.grant(ctx, ip, mac, v, Verified, now, end)
		}
		return verified, "", err
	}

	day, end := s.window(now)
	sends := 1
	prev, err := s.store.GetCode(ctx, v.Phone)
	switch {
	case errors.Is(err, ErrNotFound):
	case err != nil:
		return false, "", err
	default:
		if now.Sub(prev.SentAt) < s.cfg.ResendWait {
			return false, "", ErrTooSoon
		}
		if prev.Day == day {
			if prev.SendsToday >= s.cfg.MaxSendsPerDay {
				return false, "", ErrTooMany
			}
			sends = prev.SendsToday + 1
		}
	}

	used, err := s.store.HadFreeAccess(ctx, ip, mac, v.Phone, v.IDNumber, day)
	if err != nil {
		return false, "", err
	}
	if !used {
		if err := s.grant(ctx, ip, mac, v, Free, now, minTime(now.Add(s.cfg.FreeTime), end)); err != nil {
			return false, "", err
		}
	}

	code, err = newCode()
	if err != nil {
		return false, "", err
	}
	err = s.store.SaveCode(ctx, Code{Visitor: v, Hash: s.hash(v.Phone, code), ExpiresAt: now.Add(s.cfg.CodeTTL),
		SentAt: now, Day: day, SendsToday: sends})
	if err != nil {
		return false, "", err
	}
	return false, code, nil
}

// Verify revisa el código; si es correcto, el teléfono queda verificado y el dispositivo
// recibe internet hasta el próximo reset
func (s *Service) Verify(ctx context.Context, phone, code, ip string) error {
	ip, ok := normalizeIP(ip)
	if !phoneRe.MatchString(phone) || !codeRe.MatchString(code) || !ok {
		return ErrInvalid
	}
	now := s.now().In(s.cfg.Location)

	c, err := s.store.GetCode(ctx, phone)
	if errors.Is(err, ErrNotFound) {
		return ErrWrongCode
	}
	if err != nil {
		return err
	}
	if c.Attempts >= s.cfg.MaxAttempts {
		return ErrTooMany
	}
	if now.After(c.ExpiresAt) {
		return ErrWrongCode
	}
	if !s.codeMatches(phone, code, c.Hash) {
		if err := s.store.AddAttempt(ctx, phone); err != nil {
			return err
		}
		return ErrWrongCode
	}

	if err := s.store.SaveVerified(ctx, c.Visitor, now); err != nil {
		return err
	}
	if err := s.store.DeleteCode(ctx, phone); err != nil {
		return err
	}
	mac, err := s.mac(ctx, ip)
	if err != nil {
		return err
	}
	_, end := s.window(now)
	return s.grant(ctx, ip, mac, c.Visitor, Verified, now, end)
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

func (s *Service) grant(ctx context.Context, ip, mac string, v Visitor, kind AccessKind, now, until time.Time) error {
	day, _ := s.window(now)
	if err := s.router.GrantAccess(ctx, ip, until); err != nil {
		return err
	}
	return s.store.SaveAccess(ctx, Access{IP: ip, MAC: mac, Phone: v.Phone, IDNumber: v.IDNumber, Day: day,
		Kind: kind, Start: now, End: until})
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

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
