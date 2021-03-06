CREATE TABLE IF NOT EXISTS ra_user_passphrases (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    node_id bigint not null,
    user_id bigint not null,
    passphrase varchar(192) not null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),

    CONSTRAINT fk_ra_user_passphrases_node_id foreign key (node_id) references ra_nodes(id) on update cascade on delete cascade,
    CONSTRAINT fk_ra_user_passphrases_user_id FOREIGN KEY (user_id) REFERENCES ra_users(id) ON UPDATE cascade ON DELETE  cascade
);

CREATE INDEX IF NOT EXISTS ra_user_passphrases_user_id_index ON ra_user_passphrases USING btree(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS ra_user_passphrases_passphrase_unique_index ON ra_user_passphrases USING btree(passphrase);
CREATE INDEX IF NOT EXISTS ra_user_passphrases_inserted_at ON ra_user_passphrases USING btree(inserted_at);

CREATE TABLE IF NOT EXISTS ra_user_passphrase_invalidations (
    passphrase_id bigint PRIMARY KEY,
    node_id bigint not null,
    source_user_id bigint null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),

    CONSTRAINT fk_ra_user_passphrases_invalidations_passphrase_id FOREIGN KEY (passphrase_id)
        REFERENCES ra_user_passphrases(id) ON UPDATE cascade ON DELETE cascade,

    CONSTRAINT fk_ra_user_passphrases_invalidations_node_id foreign key (node_id)
        references ra_nodes(id) on update cascade on delete cascade,
    CONSTRAINT fk_ra_user_passphrases_invalidations_source_user_id FOREIGN KEY (source_user_id)
        REFERENCES ra_users(id) ON UPDATE cascade ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS ra_user_passphrase_invalidations_source_user_id_index ON ra_user_passphrase_invalidations USING btree(source_user_id);
