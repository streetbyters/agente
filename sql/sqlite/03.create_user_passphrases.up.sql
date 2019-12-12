CREATE TABLE ra_user_passphrases (
    id integer primary key autoincrement,
    user_id integer not null,
    passphrase text not null,
    inserted_at datetime default CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES ra_users(id) ON UPDATE cascade ON DELETE cascade
);

CREATE INDEX ra_user_passphrases_user_id_index ON ra_user_passphrases(id);
CREATE UNIQUE INDEX ra_user_passphrases_passphrase_unique_index ON ra_user_passphrases(passphrase);
