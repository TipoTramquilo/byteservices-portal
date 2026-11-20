package router

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OPNsense abre internet con la API del portal cautivo de OPNsense y lee la MAC de la tabla
// ARP. Autentica con la key y el secret de un usuario de OPNsense (Basic auth)
type OPNsense struct {
	base   string
	key    string
	secret string
	zone   int
	http   *http.Client
}

// NewOPNsense recibe la URL del equipo (https://192.168.122.10), la key y el secret de la API,
// el ID de la zona del portal, el PEM de la CA (o del certificado) del equipo si no es público
// y el nombre que el certificado debe traer en su SAN (OPNsense.internal). Sin serverName la
// verificación usa el host de la URL, así que una IP falla si el certificado no la incluye
func NewOPNsense(base, key, secret string, zone int, caPEM []byte, serverName string) (*OPNsense, error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, ServerName: serverName}}
	if len(caPEM) > 0 {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caPEM) {
			return nil, fmt.Errorf("opnsense: el PEM de la CA no tiene certificados válidos")
		}
		tr.TLSClientConfig.RootCAs = pool
	}
	return &OPNsense{
		base:   strings.TrimRight(base, "/"),
		key:    key,
		secret: secret,
		zone:   zone,
		http:   &http.Client{Transport: tr, Timeout: 10 * time.Second},
	}, nil
}

// MAC busca la IP en la tabla ARP del equipo; "" si no aparece
func (o *OPNsense) MAC(ctx context.Context, ip string) (string, error) {
	var arp []struct {
		MAC string `json:"mac"`
		IP  string `json:"ip"`
	}
	if err := o.do(ctx, http.MethodGet, "/api/diagnostics/interface/getArp", nil, &arp); err != nil {
		return "", err
	}
	for _, e := range arp {
		if e.IP == ip {
			return e.MAC, nil
		}
	}
	return "", nil
}

// GrantAccess autoriza la IP en la zona del portal. OJO: sessionTimeout todavía no está
// verificado: en la prueba con la VM la sesión quedó sin tiempo límite (acc_session_timeout
// null). Falta confirmar si la zona necesita hardtimeout o RADIUS para que venza a las 02:00
func (o *OPNsense) GrantAccess(ctx context.Context, ip string, until time.Time) error {
	secs := int(time.Until(until).Seconds())
	if secs <= 0 {
		return fmt.Errorf("opnsense: la hora de acceso para %s ya pasó", ip)
	}
	body := map[string]any{"zoneid": o.zone, "ip": ip, "sessionTimeout": secs}
	var res struct {
		ClientState string `json:"clientState"`
	}
	if err := o.do(ctx, http.MethodPost, "/api/captiveportal/session/connect", body, &res); err != nil {
		return err
	}
	if res.ClientState != "AUTHORIZED" {
		return fmt.Errorf("opnsense: sesión de %s en estado %q", ip, res.ClientState)
	}
	return nil
}

// do hace la llamada a la API; en o (si no es nil) decodifica la respuesta JSON
func (o *OPNsense) do(ctx context.Context, method, path string, in, out any) error {
	var rd io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, o.base+path, rd)
	if err != nil {
		return err
	}
	req.SetBasicAuth(o.key, o.secret)
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := o.http.Do(req)
	if err != nil {
		return fmt.Errorf("opnsense %s: %w", path, err)
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		return fmt.Errorf("opnsense %s: HTTP %d", path, res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}
