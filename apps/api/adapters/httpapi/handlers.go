// Package httpapi expone core.Service como la API JSON que llama el portal, y sirve el portal
package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"captiveportal/core"
)

type sendRequest struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Company  string `json:"company"`
	Phone    string `json:"phone"`
}

type verifyRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type response struct {
	Verified bool   `json:"verified,omitempty"`
	Redirect string `json:"redirect,omitempty"`
	Error    string `json:"error,omitempty"`
}

// New enruta la API y sirve portalDir (el dist/ compilado) desde el mismo origen, sin CORS
// redirect es el enlace de "Empezar a navegar"; vacío deja el del portal
func New(svc *core.Service, portalDir, redirect string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/send", func(w http.ResponseWriter, r *http.Request) {
		var req sendRequest
		if !decode(w, r, &req) {
			return
		}
		verified, err := svc.Send(r.Context(), core.Visitor{Name: req.Name, IDNumber: req.IDNumber, Company: req.Company, Phone: req.Phone}, clientIP(r))
		if err != nil {
			fail(w, err)
			return
		}
		res := response{Verified: verified}
		if verified {
			res.Redirect = redirect
		}
		reply(w, http.StatusOK, res)
	})

	mux.HandleFunc("POST /api/verify", func(w http.ResponseWriter, r *http.Request) {
		var req verifyRequest
		if !decode(w, r, &req) {
			return
		}
		if err := svc.Verify(r.Context(), req.Phone, req.Code, clientIP(r)); err != nil {
			fail(w, err)
			return
		}
		reply(w, http.StatusOK, response{Redirect: redirect})
	})

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })

	files := http.FileServer(http.Dir(portalDir))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { index(w, portalDir) })
	mux.Handle("GET /", files)
	return methods(mux)
}

// index sirve el portal con los endpoints de la API puestos, como lo haría la página alojada en el router
func index(w http.ResponseWriter, portalDir string) {
	page, err := os.ReadFile(filepath.Join(portalDir, "index.html"))
	if err != nil {
		log.Print(err)
		http.Error(w, "portal not built", http.StatusInternalServerError)
		return
	}
	page = bytes.Replace(page, []byte("<script>"),
		[]byte(`<script>window.BYTE_PORTAL={send:"/api/send",verify:"/api/verify"}</script><script>`), 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(page)
}

// clientIP es la IP de la conexión. No se leen X-Forwarded-For ni similares: la API recibe
// a los visitantes directo en la LAN y esas cabeceras las puede escribir cualquiera
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		reply(w, http.StatusBadRequest, response{Error: "invalid"})
		return false
	}
	return true
}

func fail(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrInvalid):
		reply(w, http.StatusBadRequest, response{Error: "invalid"})
	case errors.Is(err, core.ErrWrongCode):
		reply(w, http.StatusBadRequest, response{Error: "wrong_code"})
	case errors.Is(err, core.ErrTooSoon):
		reply(w, http.StatusTooManyRequests, response{Error: "too_soon"})
	case errors.Is(err, core.ErrTooMany):
		reply(w, http.StatusTooManyRequests, response{Error: "too_many"})
	default:
		log.Print(err)
		reply(w, http.StatusInternalServerError, response{Error: "internal"})
	}
}

func reply(w http.ResponseWriter, status int, v response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Métodos que admite cada ruta de la API; HEAD va con GET, como en el mux
var apiMethods = map[string][]string{
	"/api/send":   {http.MethodPost},
	"/api/verify": {http.MethodPost},
	"/api/health": {http.MethodGet, http.MethodHead},
}

// methods contesta 405 con Allow cuando la ruta de la API existe pero no admite ese método.
// Sin esto, el FileServer de GET / se traga los GET y HEAD a /api/* y responde 404
func methods(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, ok := apiMethods[r.URL.Path]
		if !ok || slices.Contains(allowed, r.Method) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		reply(w, http.StatusMethodNotAllowed, response{Error: "method_not_allowed"})
	})
}
