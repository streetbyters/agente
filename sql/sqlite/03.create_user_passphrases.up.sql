CREATE TABLE IF NOT EXISTS ra_user_passphrases (
    id integer primary key autoincrement,
    user_id integer not null,
    passphrase text not null,
    inserted_at datetime default CURRENT_TIMESTAMP,

    CONSTRAINT fk_ra_user_passphrases_user_id FOREIGN KEY (user_id) REFERENCES ra_users(id) ON UPDATE cascade ON DELETE cascade
);

CREATE INDEX IF NOT EXISTS ra_user_passphrases_user_id_index ON ra_user_passphrases(id);
CREATE UNIQUE INDEX IF NOT EXISTS ra_user_passphrases_passphrase_unique_index ON ra_user_passphrases(passphrase);
CREATE INDEX IF NOT EXISTS ra_user_passphrases_inserted_at ON ra_user_passphrases(inserted_at);

CREATE TABLE IF NOT EXISTS ra_user_passphrase_invalidations (
    passphrase_id integer primary key,
    source_user_id integer null,
    inserted_at datetime default CURRENT_TIMESTAMP,

    CONSTRAINT fk_ra_user_passphrases_source_user_id FOREIGN KEY (source_user_id)
        REFERENCES ra_users(id) ON UPDATE cascade ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS ra_user_passphrase_invalidations_source_user_id ON ra_user_passphrase_invalidations(source_user_id);
