-- Visitantes que pasaron el código; se saltan el código los días siguientes
CREATE TABLE IF NOT EXISTS verified (
  phone       TEXT PRIMARY KEY CHECK (phone ~ '^\+[1-9][0-9]{7,14}$'),
  name        TEXT NOT NULL CHECK (char_length(name) BETWEEN 3 AND 100),
  id_number   TEXT NOT NULL CHECK (id_number ~ '^[VE][0-9]{6,8}$'),
  company     TEXT NOT NULL CHECK (char_length(company) BETWEEN 2 AND 100),
  verified_at TIMESTAMPTZ NOT NULL
);

-- Un código pendiente por teléfono; cada envío reemplaza al anterior. Solo se guarda el hash
CREATE TABLE IF NOT EXISTS codes (
  phone       TEXT PRIMARY KEY CHECK (phone ~ '^\+[1-9][0-9]{7,14}$'),
  name        TEXT NOT NULL CHECK (char_length(name) BETWEEN 3 AND 100),
  id_number   TEXT NOT NULL CHECK (id_number ~ '^[VE][0-9]{6,8}$'),
  company     TEXT NOT NULL CHECK (char_length(company) BETWEEN 2 AND 100),
  hash        TEXT NOT NULL CHECK (char_length(hash) = 64),
  expires_at  TIMESTAMPTZ NOT NULL,
  attempts    INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0),
  sent_at     TIMESTAMPTZ NOT NULL,
  day         TEXT NOT NULL,
  sends_today INTEGER NOT NULL CHECK (sends_today >= 1)
);
