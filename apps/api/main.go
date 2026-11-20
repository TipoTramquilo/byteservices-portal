// Command captiveportal sirve el portal cautivo de Byte y su API
//
// Configuración (variables de entorno): ADDR, DB, PORTAL_DIR, TIMEZONE, RESET_AT, REDIRECT, OTP_SECRET,
// TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_WHATSAPP_FROM, TWILIO_CONTENT_SID
// Sin TWILIO_ACCOUNT_SID los códigos se imprimen en el log
package main

import (
	"crypto/rand"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	_ "time/tzdata" // los routers suelen venir sin base de zonas horarias

	"captiveportal/adapters/httpapi"
	"captiveportal/adapters/router"
	"captiveportal/adapters/sqlite"
	"captiveportal/adapters/whatsapp"
	"captiveportal/core"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC) // todas las horas en UTC, también las del log
	if len(os.Args) > 1 && os.Args[1] == "health" {
		health()
	}
	loc, err := time.LoadLocation(env("TIMEZONE", "America/Caracas"))
	if err != nil {
		log.Fatal(err)
	}
	// Hora local del reset: ahí vencen los accesos y se reinician los límites diarios
	resetAt, err := time.Parse("15:04", env("RESET_AT", "02:00"))
	if err != nil {
		log.Fatal("RESET_AT: ", err)
	}
	store, err := sqlite.Open(env("DB", "portal.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	secret := []byte(os.Getenv("OTP_SECRET"))
	if len(secret) == 0 {
		secret = make([]byte, 32)
		rand.Read(secret)
		log.Print("OTP_SECRET vacío: se usa una clave aleatoria, los códigos pendientes dejan de servir al reiniciar")
	}

	var sender core.Sender = whatsapp.Console{}
	if sid := os.Getenv("TWILIO_ACCOUNT_SID"); sid != "" {
		sender = whatsapp.NewTwilio(sid, os.Getenv("TWILIO_AUTH_TOKEN"),
			os.Getenv("TWILIO_WHATSAPP_FROM"), os.Getenv("TWILIO_CONTENT_SID"))
	}

	svc := core.NewService(store, sender, router.Log{}, core.Config{
		Location:       loc,
		ResetHour:      resetAt.Hour(),
		ResetMinute:    resetAt.Minute(),
		FreeTime:       time.Hour,
		CodeTTL:        10 * time.Minute,
		ResendWait:     60 * time.Second,
		MaxSendsPerDay: 5,
		MaxAttempts:    5,
		Secret:         secret,
	})

	srv := &http.Server{
		Addr:              env("ADDR", ":8080"),
		Handler:           httpapi.New(svc, env("PORTAL_DIR", "../web/dist"), os.Getenv("REDIRECT")),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	log.Printf("listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// health es el healthcheck del contenedor: la imagen distroless no tiene shell ni curl
func health() {
	_, port, _ := net.SplitHostPort(env("ADDR", ":8080"))
	res, err := (&http.Client{Timeout: 3 * time.Second}).Get("http://127.0.0.1:" + port + "/api/health")
	if err != nil || res.StatusCode/100 != 2 {
		os.Exit(1)
	}
	os.Exit(0)
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
