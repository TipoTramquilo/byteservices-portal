package core

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func newCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// hash ata el código a su teléfono, así un hash filtrado no sirve para otro número
func (s *Service) hash(phone, code string) string {
	m := hmac.New(sha256.New, s.cfg.Secret)
	m.Write([]byte(phone + ":" + code))
	return hex.EncodeToString(m.Sum(nil))
}

func (s *Service) codeMatches(phone, code, hash string) bool {
	return hmac.Equal([]byte(s.hash(phone, code)), []byte(hash))
}
