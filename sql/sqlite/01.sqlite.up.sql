CREATE TABLE IF NOT EXISTS ra_migrations (
  id integer primary key autoincrement,
  number integer not null,
  name text not null,
  inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
