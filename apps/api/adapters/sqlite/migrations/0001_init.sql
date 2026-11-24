-- Tablas STRICT: SQLite rechaza valores del tipo equivocado en vez de convertirlos
-- Fechas y horas: ISO 8601 en UTC ('2026-11-19T18:03:00Z'); días: ISO 8601 ('2026-11-19')
-- en la zona horaria del portal. datetime()/date() devuelven NULL si el texto no es válido
-- Teléfonos: E.164 ('+' y hasta 15 dígitos, sin 0 inicial). IP: IPv4 o IPv6 en forma canónica (la valida la API)
-- MAC: EUI-48 en mayúsculas con ':'

-- Visitantes que enviaron el formulario; el teléfono queda registrado para el evento
CREATE TABLE IF NOT EXISTS verified (
  phone       TEXT PRIMARY KEY CHECK (phone GLOB '+[1-9]*' AND length(phone) BETWEEN 9 AND 16
                                      AND substr(phone, 2) NOT GLOB '*[^0-9]*'),
  name        TEXT NOT NULL CHECK (length(name) BETWEEN 3 AND 100),
  id_number   TEXT NOT NULL CHECK (id_number GLOB '[VE]*' AND length(id_number) BETWEEN 7 AND 9
                                   AND substr(id_number, 2) NOT GLOB '*[^0-9]*'),  -- V18456789
  company     TEXT NOT NULL CHECK (length(company) BETWEEN 2 AND 100),
  verified_at TEXT NOT NULL CHECK (datetime(verified_at) IS NOT NULL)
) STRICT;

-- Cada vez que se da internet (una sesión de 30 min al enviar el formulario): registro del evento
CREATE TABLE IF NOT EXISTS access (
  id        INTEGER PRIMARY KEY,
  ip        TEXT NOT NULL CHECK (length(ip) BETWEEN 3 AND 39 AND ip NOT GLOB '*[^0-9a-f.:]*'),
  -- MAC: EUI-48 en mayúsculas con ':'. La consulta el backend al router; NULL si no la conoce
  mac       TEXT CHECK (mac IS NULL OR mac GLOB '[0-9A-F][0-9A-F]:[0-9A-F][0-9A-F]:[0-9A-F][0-9A-F]:[0-9A-F][0-9A-F]:[0-9A-F][0-9A-F]:[0-9A-F][0-9A-F]'),
  phone     TEXT NOT NULL CHECK (phone GLOB '+[1-9]*' AND length(phone) BETWEEN 9 AND 16
                                 AND substr(phone, 2) NOT GLOB '*[^0-9]*'),
  id_number TEXT NOT NULL CHECK (id_number GLOB '[VE]*' AND length(id_number) BETWEEN 7 AND 9
                                 AND substr(id_number, 2) NOT GLOB '*[^0-9]*'),
  day       TEXT NOT NULL CHECK (date(day) = day),
  kind      TEXT NOT NULL CHECK (kind = 'verified'),
  start_at  TEXT NOT NULL CHECK (datetime(start_at) IS NOT NULL),
  end_at    TEXT NOT NULL CHECK (datetime(end_at) IS NOT NULL AND end_at > start_at)
) STRICT;
