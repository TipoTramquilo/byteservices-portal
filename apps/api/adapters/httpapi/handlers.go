// Package httpapi expone core.Service como la API JSON que llama el portal, y sirve el portal
package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"html"
	"log"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"byte-api/core"
)

type sendRequest struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Company  string `json:"company"`
	Phone    string `json:"phone"`
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
		ip := clientIP(r)
		// ip ya pasó por netip.ParseAddr y las cabeceras van con %q (escapadas): no hay inyección en el log
		log.Printf("send desde %s (peer=%s X-Real-Ip=%q X-Forwarded-For=%q)", ip, r.RemoteAddr, r.Header.Get("X-Real-Ip"), r.Header.Get("X-Forwarded-For")) //nolint:gosec // G706: ver arriba
		if err := svc.Send(r.Context(), core.Visitor{Name: req.Name, IDNumber: req.IDNumber, Company: req.Company, Phone: req.Phone}, ip); err != nil {
			fail(w, err)
			return
		}
		// Sin código: la sesión ya está abierta, así que el portal pasa directo a "conectado"
		reply(w, http.StatusOK, response{Verified: true, Redirect: redirect})
	})

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })

	files := http.FileServer(http.Dir(portalDir))
	network := os.Getenv("NETWORK_NAME")
	if network == "" {
		network = "Byte Invitados"
	}
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { index(w, portalDir, network) })
	// el redirect del portal cautivo llega a /index.html: también lleva la configuración del portal
	mux.HandleFunc("GET /index.html", func(w http.ResponseWriter, r *http.Request) { index(w, portalDir, network) })
	mux.Handle("GET /", files)
	return cors(captive(methods(mux), os.Getenv("CAPTIVE_HOME")), os.Getenv("CORS_ORIGIN"))
}

// index sirve el portal con los endpoints de la API puestos, como lo haría la página alojada en el router.
// network reemplaza el nombre de la red ("Byte Invitados" en la página) en el título y la pastilla
func index(w http.ResponseWriter, portalDir, network string) {
	page, err := os.ReadFile(filepath.Join(portalDir, "index.html")) //nolint:gosec // G304: PORTAL_DIR lo fija el operador, no el visitante
	if err != nil {
		log.Print(err)
		http.Error(w, "portal not built", http.StatusInternalServerError)
		return
	}
	page = bytes.ReplaceAll(page, []byte("Byte Invitados"), []byte(html.EscapeString(network)))
	page = bytes.Replace(page, []byte("<script>"),
		[]byte(`<script>window.BYTE_PORTAL={send:"/api/send"}</script><script>`), 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(page) // si el visitante cerró la conexión no hay a quién avisarle
}

// clientIP es la IP del visitante. Normalmente es la de la conexión. Si la conexión viene de
// 127.0.0.1 (el lighttpd del portal reenvía a la API en local), la IP real sale de X-Real-Ip o
// de la primera IP de X-Forwarded-For. Desde otras IP esas cabeceras se ignoran: las puede escribir cualquiera
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if ip, perr := netip.ParseAddr(host); perr == nil && ip.IsLoopback() {
		if real, ok := forwardedIP(r); ok {
			return real
		}
	}
	return host
}

// forwardedIP lee la IP real de las cabeceras que pone el proxy local, si es una IP válida
func forwardedIP(r *http.Request) (string, bool) {
	candidates := []string{r.Header.Get("X-Real-Ip")}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		candidates = append(candidates, strings.TrimSpace(strings.Split(xff, ",")[0]))
	}
	for _, c := range candidates {
		if ip, err := netip.ParseAddr(strings.TrimSpace(c)); err == nil {
			return ip.String(), true
		}
	}
	return "", false
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
	default:
		log.Print(err)
		reply(w, http.StatusInternalServerError, response{Error: "internal"})
	}
}

func reply(w http.ResponseWriter, status int, v response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) // el status ya salió: un error aquí es la conexión cortada
}

// Métodos que admite cada ruta de la API; HEAD va con GET, como en el mux
var apiMethods = map[string][]string{
	"/api/send":   {http.MethodPost},
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

// cors deja que la página alojada en el router (otro origen) llame a la API. Solo acepta el
// origen de CORS_ORIGIN; sin la variable no se agregan encabezados
func cors(next http.Handler, origin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin != "" && r.Header.Get("Origin") == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Vary", "Origin")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// captive manda al portal a quien pide otra cosa: el router redirige todo el tráfico sin
// autorizar aquí, y el SO del teléfono solo abre el portal si recibe una redirección (302)
func captive(next http.Handler, home string) http.Handler {
	if home == "" {
		return next
	}
	u, err := url.Parse(home)
	if err != nil {
		log.Fatal("CAPTIVE_HOME: ", err)
	}
	portal := u.Hostname()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if i := strings.LastIndex(host, ":"); i > 0 {
			host = host[:i]
		}
		if host == portal || host == "127.0.0.1" || host == "localhost" {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, home+"/", http.StatusFound)
	})
}
