// Command captiveportal sirve el portal cautivo de Byte y su API
//
// Configuración (variables de entorno): ADDR, DB, PORTAL_DIR, TIMEZONE, RESET_AT, REDIRECT,
// OPNSENSE_URL, OPNSENSE_KEY, OPNSENSE_SECRET, OPNSENSE_ZONE, OPNSENSE_CA_FILE, OPNSENSE_SERVER_NAME
// Sin OPNSENSE_URL el acceso no se abre (solo se registra en el log)
package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	_ "time/tzdata" // los routers suelen venir sin base de zonas horarias

	"captiveportal/adapters/httpapi"
	"captiveportal/adapters/router"
	"captiveportal/adapters/sqlite"
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

	var rtr core.Router = router.Log{}
	if url := os.Getenv("OPNSENSE_URL"); url != "" {
		var ca []byte
		if f := os.Getenv("OPNSENSE_CA_FILE"); f != "" {
			if ca, err = os.ReadFile(f); err != nil {
				log.Fatal("OPNSENSE_CA_FILE: ", err)
			}
		}
		zone, err := strconv.Atoi(env("OPNSENSE_ZONE", "0"))
		if err != nil {
			log.Fatal("OPNSENSE_ZONE: ", err)
		}
		rtr, err = router.NewOPNsense(url, os.Getenv("OPNSENSE_KEY"), os.Getenv("OPNSENSE_SECRET"), zone, ca,
			os.Getenv("OPNSENSE_SERVER_NAME"))
		if err != nil {
			log.Fatal(err)
		}
	}

	svc := core.NewService(store, rtr, core.Config{
		Location:    loc,
		ResetHour:   resetAt.Hour(),
		ResetMinute: resetAt.Minute(),
		SessionTime: 30 * time.Minute,
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
