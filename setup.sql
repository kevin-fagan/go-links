CREATE TABLE IF NOT EXISTS links (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    short_url       TEXT NOT NULL UNIQUE,
    long_url        TEXT NOT NULL,
    visits          INTEGER DEFAULT 0,
    last_updated    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS logs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    short_url   TEXT, 
    long_url    TEXT,
    tag         TEXT,
    client_ip   TEXT NOT NULL,
    action      TEXT NOT NULL CHECK(action IN ('CREATE', 'UPDATE', 'DELETE')),
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS links_tags (
    tag_id INTEGER,
    link_id INTEGER,
    PRIMARY KEY (tag_id, link_id),
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);
