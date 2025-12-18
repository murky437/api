CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW')),
	updated_at TEXT
, is_guest BOOLEAN NOT NULL DEFAULT 0);
CREATE TABLE refresh_token (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER REFERENCES user(id) ON DELETE CASCADE,
	jwt TEXT NOT NULL,
	expires_at TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW'))
);
CREATE TABLE project (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER REFERENCES user(id) ON DELETE CASCADE,
	title TEXT NOT NULL,
	slug TEXT NOT NULL UNIQUE,
	notes TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')),
	updated_at TEXT
, sort_index INTEGER NOT NULL DEFAULT 0);
CREATE INDEX idx_project_user_id ON project(user_id);
CREATE TABLE version (
	number INTEGER PRIMARY KEY,
	applied_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ', 'NOW'))
);
CREATE TABLE long_reminder (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	interval_days INTEGER NOT NULL,
	user_id INTEGER NOT NULL REFERENCES user(id),
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')),
	updated_at TEXT,
	last_reminded_at TEXT,
	marked_done_at TEXT,
	is_enabled BOOLEAN NOT NULL DEFAULT 1
);
CREATE INDEX idx_long_reminder_user_id ON long_reminder(user_id);
CREATE TABLE guest_login_token (
	token TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW'))
);
CREATE TABLE daily_brevo_sends (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    number INTEGER NOT NULL,
    date TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL DEFAULT(STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')),
    updated_at TEXT NOT NULL
);
