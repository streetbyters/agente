CREATE TABLE  ra_users (
    id integer primary key autoincrement,
    username text not null,
    password_digest text not null,
    email text not null,
    is_active numeric default 1,
    inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX  ra_users_username_unique_index ON ra_users(username);
CREATE UNIQUE INDEX  ra_users_email_unique_index ON ra_users(email);

CREATE TABLE  ra_jobs (
    id integer primary key autoincrement,
    source_user_id integer null,
    inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (source_user_id) REFERENCES ra_users(id) ON UPDATE cascade ON DELETE SET NULL
);

CREATE INDEX  ra_jobs_source_user_id_index ON ra_jobs(source_user_id);

CREATE TABLE  ra_job_details (
    id integer primary key autoincrement,
    job_id integer not null,
    source_user_id integer null,

    name text not null,
    type text default 'new_release',
    detail text null,
    before numeric default 0,
    after numeric default 0,

    script_file text null,
    script text null,

    inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (job_id) REFERENCES ra_jobs(id) ON UPDATE cascade  ON DELETE cascade,
    FOREIGN KEY (source_user_id) REFERENCES ra_users(id) ON UPDATE cascade ON DELETE SET NULL
);

CREATE INDEX  ra_job_details_job_id_index ON ra_job_details(job_id);
CREATE INDEX  ra_job_details_source_user_id_index ON ra_job_details(source_user_id);

CREATE TABLE  ra_job_logs(
  id integer primary key autoincrement,
  job_id integer not null,
  data text null,
  state numeric default 0,
  inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (job_id) REFERENCES ra_jobs(id) ON UPDATE cascade ON DELETE CASCADE
);

CREATE INDEX  ra_job_logs_job_id_index ON ra_job_logs(job_id);
CREATE INDEX  ra_job_logs_inserted_at_index ON ra_job_logs(inserted_at DESC);
