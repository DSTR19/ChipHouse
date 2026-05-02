CREATE TABLE IF NOT EXISTS studios (
    id      SERIAL PRIMARY KEY,
    number  TEXT    NOT NULL UNIQUE,
    area    NUMERIC NOT NULL CHECK (area > 0),
    rent    NUMERIC NOT NULL CHECK (rent >= 0)
);

CREATE TABLE IF NOT EXISTS tenants (
    id          SERIAL    PRIMARY KEY,
    full_name   TEXT      NOT NULL,
    studio_id   INT       NOT NULL REFERENCES studios(id) ON DELETE RESTRICT,
    move_in     DATE      NOT NULL,
    move_out    DATE,
    persons     INT       NOT NULL DEFAULT 1 CHECK (persons >= 1),
    is_active   BOOLEAN   NOT NULL DEFAULT TRUE
);

CREATE UNIQUE INDEX IF NOT EXISTS tenants_one_active_per_studio
    ON tenants (studio_id) WHERE is_active = TRUE;

CREATE TABLE IF NOT EXISTS users (
    id            SERIAL       PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    role          TEXT         NOT NULL CHECK (role IN ('USER', 'ADMIN', 'INVESTOR')),
    tenant_id     INT          REFERENCES tenants(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS charges (
    id           SERIAL    PRIMARY KEY,
    studio_id    INT       NOT NULL REFERENCES studios(id) ON DELETE RESTRICT,
    tenant_id    INT       NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    period_start DATE      NOT NULL,
    period_end   DATE      NOT NULL,
    rent         NUMERIC   NOT NULL DEFAULT 0,
    utilities    NUMERIC   NOT NULL DEFAULT 0,
    adjustment   NUMERIC   NOT NULL DEFAULT 0,
    total        NUMERIC   NOT NULL DEFAULT 0,
    paid         NUMERIC   NOT NULL DEFAULT 0,
    issued_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, period_start)
);

CREATE TABLE IF NOT EXISTS complaints (
    id          SERIAL      PRIMARY KEY,
    tenant_id   INT         NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    studio_id   INT         NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
    subject     TEXT        NOT NULL,
    body        TEXT        NOT NULL,
    status      TEXT        NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'RESOLVED')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS vacancy_periods (
    id         SERIAL  PRIMARY KEY,
    studio_id  INT     NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
    start_date DATE    NOT NULL,
    end_date   DATE,
    lost_rent  NUMERIC NOT NULL DEFAULT 0
);
