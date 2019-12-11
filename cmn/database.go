// Copyright 2019 Abdulkadir DILSIZ - TransferChain
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmn

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/akdilsiz/agente/model"
	_ "github.com/go-sql-driver/mysql" // Mysql Driver
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // Postgres Driver
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // SQLite Driver
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// Tables db table enum type
type Tables string

// DB Table enums
const (
	tUser      Tables = "ra_users"
	tJob       Tables = "ra_jobs"
	tJobDetail Tables = "ra_job_details"
	tJobLog    Tables = "ra_job_logs"
	tMigration Tables = "ra_migrations"
)

// Error for database violation errors
type Error int

const (
	// TableNotFound sql violation code
	TableNotFound Error = 1
	// OtherError unhandled sql violation codes
	OtherError Error = 0
	// InternalError SQLite Error
	InternalError Error = 1
)

// Database struct
type Database struct {
	Config *model.Config
	Type   model.DB
	DB     *sqlx.DB
	Tx     *sqlx.Tx
	Error  error
}

// Tx transaction for database queries
type Tx struct {
	DB *Database
}

// Result structure for database query results
type Result struct {
	Rows  []interface{}
	Count int64
	Error error
}

// NewDB building database
func NewDB(config *model.Config, logger *Logger) (*Database, error) {
	database := &Database{}
	database.Config = config

	switch config.DB {
	case model.SQLite:
		connURL := fmt.Sprintf("file:%s?cache=shared&mode=rwc",
			filepath.Join(config.DBPath, config.DBName))

		var db *sqlx.DB
		var db2 *sql.DB

		db2, _ = sql.Open("sqlite3", connURL)
		db = sqlx.NewDb(db2, "sqlite3")

		if err := db.Ping(); err != nil {
			return nil, err
		}

		database.Type = model.SQLite
		database.DB = db
		break
	case model.Postgres:
		connURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.DBHost,
			config.DBPort,
			config.DBUser,
			config.DBPass,
			config.DBName,
			config.DBSsl)

		db, _ := sqlx.Open("postgres", connURL)

		if err := db.Ping(); err != nil {
			return nil, err
		}

		db.SetMaxIdleConns(15)
		db.SetMaxOpenConns(15)

		database.Type = model.Postgres
		database.DB = db
		break
	case model.Mysql:
		connURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=%s",
			config.DBUser,
			config.DBPass,
			config.DBHost,
			config.DBPort,
			config.DBName,
			config.DBSsl)

		db, _ := sqlx.Open("mysql", connURL)

		if err := db.Ping(); err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(15)
		db.SetMaxOpenConns(15)

		database.Type = model.Mysql
		database.DB = db

		break
	default:
		return nil, errors.New("unsupported database")
	}

	return database, nil
}

// DropDB Drop database schemas
func DropDB(database *Database) error {
	var err error
	switch database.Type {
	case model.SQLite, model.Postgres, model.Mysql:
		files := migrationFiles(database, "down")
		for _, f := range files {
			result := database.Query(f.Data)
			err = result.Error
		}
		break
	}
	return err
}

// InstallDB Database schemas installer
func InstallDB(database *Database) error {
	var err error

	switch database.Type {
	case model.SQLite, model.Postgres, model.Mysql:
		err = migrationUp(database)
		break
	}

	return err
}

type sqlS struct {
	Number int
	Name   string
	Data   string
}

func migrationFiles(db *Database, typ string) []sqlS {
	var sqls []sqlS

	var files []string

	files, _ = filepath.Glob(filepath.Join(db.Config.DBPath, "sql", string(db.Config.DB), "[0-9]*.*."+typ+".sql"))

	for _, f := range files {
		fileName := strings.Split(f, "/")[len(strings.Split(f, "/"))-1]
		fileNumber := strings.Split(fileName, ".")[0]
		n, _ := strconv.Atoi(fileNumber)
		data, err := ioutil.ReadFile(f)
		if err == nil {
			sqls = append(sqls, sqlS{
				Number: n,
				Name:   fileName,
				Data:   string(data),
			})
		}
	}

	return sqls
}

func migrationUp(db *Database) error {
	var err error
	result := Result{}
	files := migrationFiles(db, "up")

	_, err = db.DB.Queryx("SELECT * FROM " + string(tMigration) + " AS m ORDER BY id ASC")
	if err != nil {
		if  int(dbError(db, err)) == int(TableNotFound) {
			err = nil
			for _, f := range files {
				switch f.Name {
				case "01.postgres.up.sql", "01.sqlite.up.sql", "01.mysql.up.sql":
					err = db.DB.QueryRowx(f.Data).Err()
					break
				}
			}

			tx, _ := db.DB.Beginx()

			for _, f := range files {
				switch f.Name {
				case "01.postgres.up.sql", "01.sqlite.up.sql", "01.mysql.up.sql":
					break
				default:
					_, err = tx.Queryx(f.Data)
					if err != nil {
						err = tx.Rollback()
						break
					}

					err = tx.QueryRowx("INSERT INTO " + string(tMigration) + " (" +
						"number, name) VALUES ($1, $2)", f.Number, f.Name).Err()
				}
			}

			if err != nil {
				return err
			}
			tx.Commit()
		}

		return nil
	}

	err = nil

	result = db.Query("SELECT * FROM " + string(tMigration) + " AS m ORDER BY id ASC")
	var lastMigration []interface{}
	if len(result.Rows) > 0 {
		lastMigration = result.Rows[:len(result.Rows)]
	}

	tx, err := db.DB.Beginx()
	files = migrationFiles(db, "up")

	for _, f := range files {
		switch f.Name {
		case "01.postgres.up.sql", "01.sqlite.up.sql", "01.mysql.up.sql":

			break
		default:
			if len(lastMigration) > 0 {
				if f.Number > int(lastMigration[1].(int64)) {
					err = tx.QueryRowx(f.Data).Err()
					if err != nil {
						tx.Rollback()
						break
					}

					err = tx.QueryRowx("INSERT INTO " + string(tMigration) + " (" +
						"number, name) VALUES ($1, $2)", f.Number, f.Name).Err()
				}
			} else {
				err = tx.QueryRowx(f.Data).Err()
				if err != nil {
					tx.Rollback()
					break
				}

				err = tx.QueryRowx("INSERT INTO " + string(tMigration) + " (" +
					"number, name) VALUES ($1, $2)", f.Number, f.Name).Err()
			}
		}
	}

	tx.Commit()

	return err
}

func dbError(db *Database, err error) Error {
	switch db.Type {
	case model.Postgres:
		if pgerr, ok := err.(*pq.Error); ok {
			switch string(pgerr.Code) {
			case "42P01":
				return TableNotFound
			default:
				return OtherError
			}
		}
		break
	case model.Mysql:

		break
	case model.SQLite:
		if SQLiteErr, ok := err.(sqlite3.Error); ok {
			switch SQLiteErr.Code.Error() {
			case "SQL logic error":
				return InternalError
			}
		}
		break
	}

	return -1
}

func (d *Database) beginTx() *Database {
	if d.Tx == nil {
		tx, err := d.DB.Beginx()
		if err != nil {
			d.Error = err
		}
		d.Tx = tx
		return d
	}
	d.Error = nil
	return d
}

func (d *Database) rollback() *Database {
	if d.Tx != nil {
		if err := d.Tx.Rollback(); err != nil {
			d.Error = err
			return d
		}
	}
	d.Error = nil
	return d
}

func (d *Database) commit() *Database {
	if d.Tx != nil {
		if err := d.Tx.Commit(); err != nil {
			d.Error = err
			return d
		}
		d.Tx = nil
		d.Error = nil
	}
	return d
}

// Query database query builder
func (d *Database) Query(query string, params ...interface{}) Result {
	result := Result{}

	if d.Error != nil {
		result.Error = d.Error
		return result
	}

	var rows *sqlx.Rows
	var err error

	if d.Tx != nil {
		rows, err = d.Tx.Queryx(query, params...)
	} else {
		rows, err = d.DB.Queryx(query, params...)
	}

	if err != nil {
		d.rollback()
		d.Error = err
		result.Error = err
		return result
	}

	for rows.Next() {
		result.Rows, result.Error = rows.SliceScan()
	}
	defer rows.Close()

	return result
}

// QueryRow database row query builder
func (d *Database) QueryRow(query string, params ...interface{}) Result {
	result := Result{}

	var r interface{}
	var err error
	if d.Tx != nil {
		err = d.Tx.QueryRowx(query, params...).Scan(&r)
	} else {
		err = d.DB.QueryRowx(query, params...).Scan(&r)
	}
	if err != nil {
		d.rollback()
		result.Error = err
	}
	result.Rows = append(result.Rows, r)

	return result
}

// Transaction database tx builder
func (d *Database) Transaction(cb func(tx *Tx) error) *Database {
	d.beginTx()
	newTx := new(Tx)
	newTx.DB = d
	if cb(newTx) != nil {
		return d.rollback()
	}
	return d.commit()
}

// Select query builder by database type.
func (t *Tx) Select(table Tables, whereClause string) Result {
	result := Result{}

	return result
}

// Insert query builder by database type
func (t *Tx) Insert(table Tables, data interface{}) Result {
	result := Result{}

	return result
}

// Update query builder by database type
func (t *Tx) Update(table Tables, whereClause string, data interface{}) Result {
	result := Result{}

	return result
}

// Delete query build by database type
func (t *Tx) Delete(table Tables, whereClause string) Result {
	result := Result{}

	return result
}
