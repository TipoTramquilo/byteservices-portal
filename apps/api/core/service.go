package core

import (
	"context"
	"errors"
	"hash/fnv"
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
	Location       *time.Location // zona horaria del día; todo se guarda en UTC
	ResetHour      int            // hora local en que empieza un día nuevo (2 = 2:00 AM)
	ResetMinute    int
	CodeTTL        time.Duration
	ResendWait     time.Duration // espera mínima entre códigos del mismo teléfono
	MaxSendsPerDay int
	MaxAttempts    int
	Secret         []byte // clave HMAC de los hashes de código
}

type Service struct {
	store  Store
	sender Sender
	cfg    Config
	now    func() time.Time
	locks  [phoneLocks]sync.Mutex
}

// phoneLocks son las franjas del lock por teléfono: dos números pueden compartir franja, lo que
// solo cuesta algo de espera entre ellos
const phoneLocks = 256

// lockPhone bloquea ese teléfono y devuelve la función que lo suelta
func (s *Service) lockPhone(phone string) func() {
	h := fnv.New32a()
	h.Write([]byte(phone))
	m := &s.locks[h.Sum32()%phoneLocks]
	m.Lock()
	return m.Unlock
}

func NewService(store Store, sender Sender, cfg Config) *Service {
	return &Service{store: store, sender: sender, cfg: cfg, now: time.Now}
}

// Los mismos formatos que exige la base
var (
	idNumberRe = regexp.MustCompile(`^[VE]\d{6,8}$`)
	phoneRe    = regexp.MustCompile(`^\+[1-9]\d{7,14}$`) // E.164
	codeRe     = regexp.MustCompile(`^\d{6}$`)
)

// Send valida el formulario y manda un código por WhatsApp. Un teléfono ya verificado no recibe
// código: responde verified=true y no manda nada
func (s *Service) Send(ctx context.Context, v Visitor) (verified bool, err error) {
	v.Name = strings.Join(strings.Fields(v.Name), " ")
	v.Company = strings.Join(strings.Fields(v.Company), " ")
	if len(strings.Fields(v.Name)) < 2 || utf8.RuneCountInString(v.Name) > 100 || !idNumberRe.MatchString(v.IDNumber) ||
		utf8.RuneCountInString(v.Company) < 2 || utf8.RuneCountInString(v.Company) > 100 ||
		!phoneRe.MatchString(v.Phone) {
		return false, ErrInvalid
	}
	now := s.now().In(s.cfg.Location)

	// el chequeo y el guardado del código van juntos: si dos peticiones del mismo teléfono se
	// cruzan, ambas ven el código anterior y cada una manda el suyo. El lock se suelta antes de
	// llamar al proveedor, que es lo lento y no toca el estado
	unlock := s.lockPhone(v.Phone)
	verified, code, err := s.prepare(ctx, v, now)
	unlock()
	if err != nil || verified {
		return verified, err
	}
	return false, s.sender.SendCode(ctx, v.Phone, code)
}

// prepare decide si toca enviar y guarda el código; devuelve el código a mandar (vacío si no hay
// que mandar nada). Se llama con el lock del teléfono tomado
func (s *Service) prepare(ctx context.Context, v Visitor, now time.Time) (verified bool, code string, err error) {
	verified, err = s.store.IsVerified(ctx, v.Phone)
	if err != nil || verified {
		return verified, "", err
	}

	day, _ := s.window(now)
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

// Verify revisa el código; si es correcto, el teléfono queda verificado y se guarda el visitante
func (s *Service) Verify(ctx context.Context, phone, code string) error {
	if !phoneRe.MatchString(phone) || !codeRe.MatchString(code) {
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
	return s.store.DeleteCode(ctx, phone)
}

// window devuelve el día que contiene now: va de un reset al siguiente. day es la fecha local
// en que empezó (a la 1:00 AM del 21 sigue siendo el día "2026-11-20")
func (s *Service) window(now time.Time) (day string, end time.Time) {
	y, m, d := now.Date()
	start := time.Date(y, m, d, s.cfg.ResetHour, s.cfg.ResetMinute, 0, 0, s.cfg.Location)
	if now.Before(start) {
		start = time.Date(y, m, d-1, s.cfg.ResetHour, s.cfg.ResetMinute, 0, 0, s.cfg.Location)
	}
	sy, sm, sd := start.Date()
	return start.Format(time.DateOnly), time.Date(sy, sm, sd+1, s.cfg.ResetHour, s.cfg.ResetMinute, 0, 0, s.cfg.Location)
}
