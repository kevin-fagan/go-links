CREATE TABLE IF NOT EXISTS links (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    short_url       TEXT NOT NULL UNIQUE,
    long_url        TEXT NOT NULL,
    visits          INTEGER DEFAULT 0,
    last_updated    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS audit (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    short_url   TEXT NOT NULL, 
    long_url    TEXT NOT NULL,
    action      TEXT NOT NULL CHECK(action IN ('CREATE', 'UPDATE', 'DELETE')),
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP
);

