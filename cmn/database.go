//
// Copyright 2019 Abdulkadir DILSIZ
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

// Package cmn provides in-app modules
package cmn

import (
	"errors"
	"fmt"
	"github.com/akdilsiz/release-agent/model"
	_ "github.com/go-sql-driver/mysql" // Mysql Driver
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // Postgres Driver
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Tables db table enum type
type Tables string

// DB Table enums
const tUser Tables = "ra_users"
const tJob Tables = "ra_jobs"
const tJobDetail Tables = "ra_job_details"
const tJobLog Tables = "ra_job_logs"
const tMigration Tables = "ra_migrations"

// Database struct
type Database struct {
	Config		*model.Config
	Type 		model.DB
	Bolt		*bolt.DB
	DB			*sqlx.DB
	Tx			*sqlx.Tx
	Error		error
}

// Result structure for database query results
type Result struct {
	Rows 	[]interface{}
	Count	int64
	Error	error
}

// NewDB Database connection with config struct and logger package
func NewDB(config *model.Config, logger *Logger) (*Database, error) {
	database := &Database{}
	database.Config = config

	switch config.DB {
	case model.Bolt:
		db, err := bolt.Open(path.Join(config.DBPath, config.DBName), 0666, nil)
		if err != nil {
			return nil, err
		}
		database.Type = model.Bolt
		database.Bolt = db
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
	case model.Bolt:
		database.Bolt.Update(func(tx *bolt.Tx) error {
			tx.DeleteBucket([]byte(tMigration))
			tx.DeleteBucket([]byte(tUser))
			tx.DeleteBucket([]byte(tJob))
			tx.DeleteBucket([]byte(tJobDetail))
			tx.DeleteBucket([]byte(tJobLog))
			return nil
		})
		break
	case model.Postgres, model.Mysql:
		files := migrationFiles(database, "down")
		for _, f := range files {
			result := database.Query(f)
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
	case model.Bolt:
		err = database.Bolt.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(tMigration))
			if b == nil {
				_, err := tx.CreateBucket([]byte(tUser))

				if err != nil {
					return err
				}
				_, err = tx.CreateBucket([]byte(tJob))
				if err != nil {
					return err
				}
				_, err = tx.CreateBucket([]byte(tJobDetail))
				if err != nil {
					return err
				}
				_, err = tx.CreateBucket([]byte(tJobLog))
				if err != nil {
					return err
				}
			}

			_, err = tx.CreateBucket([]byte(tMigration))
			return err
		})
		break
	case model.Postgres, model.Mysql:
		result := database.Query("SELECT * FROM " + string(tMigration) + " AS m ORDER BY id ASC")
		if result.Error != nil {
			if dbError(database, result.Error) == 1 {
				files := migrationFiles(database, "up")
				for _, f := range files {
					database := database.BeginTx()
					result := database.Query(f)
					if result.Error != nil {
						database.Rollback()
						break
					}
					database.Commit()
				}
				return nil
			} else {
				return result.Error
			}
		} else {
			if len(result.Rows) > 0 {
				//lastMigration := result.Rows[:len(result.Rows)]
			} else {
				files := migrationFiles(database, "up")
				for _, f := range files {
					database := database.BeginTx()
					result = database.Query(f)
					if result.Error != nil {
						database.Rollback()
						break
					}
					database.Commit()
				}
				return nil
			}
		}
		break
	default:
		break
	}

	return err
}

func migrationFiles(db *Database, typ string) map[int]string {
	sqls := make(map[int]string)

	var files []string
	files, _ = filepath.Glob(filepath.Join(db.Config.DBPath, string(db.Config.DB), "[0-9]*.*."+typ+".sql"))
	for _, f := range files {
		fileName := strings.Split(f, "/")[len(strings.Split(f, "/")) - 1]
		fileNumber := strings.Split(fileName, ".")[0]
		n, _ := strconv.Atoi(fileNumber)
		data, err := ioutil.ReadFile(f)
		if err == nil {
			sqls[n] = string(data)
		}
	}
	return sqls
}

func dbError(db *Database, err error) int {
	switch db.Type {
	case model.Postgres:
		if pgerr, ok := err.(*pq.Error); ok {
			switch string(pgerr.Code) {
			case "42P01":
				return 1
			default:
				return 0
			}
		}
		break
	case model.Mysql:

		break
	}

	return -1
}

func (d *Database) BeginTx() *Database {
	if d.Tx == nil {
		tx, err := d.DB.Beginx()
		if err != nil {
			d.Error = err
		}
		d.Tx = tx
	}
	return d
}

func (d *Database) Rollback() *Database {
	if d.Tx != nil {
		if err := d.Tx.Rollback(); err != nil {
			d.Error = err
			return d
		}
	}
	d.Error = nil
	return d
}

func (d *Database) Commit() *Database {
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

// Query Raw query builder. Returns multiple records
func (d *Database) Query(query string, params ...interface{}) Result {
	result := Result{}

	if d.Type == model.Bolt {
		result.Error = errors.New("this method not supported")
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
		result.Error = err
		return result
	}
	for rows.Next() {
		result.Rows, result.Error = rows.SliceScan()
	}
	defer rows.Close()

	return result
}

// QueryRow Raw query builder. Returns one record.
func (d *Database) QueryRow(query string, params ...interface{}) Result {
	result := Result{}

	if d.Type == model.Bolt {
		result.Error = errors.New("this method not supported")
		return result
	}

	var r interface{}
	var err error
	if d.Tx != nil {
		err = d.Tx.QueryRowx(query, params...).Scan(&r)
	} else {
		err = d.DB.QueryRowx(query, params...).Scan(&r)
	}
	if err != nil {
		d.Rollback()
		result.Error = err
	}
	result.Rows[0] = r

	return result
}

// Select query builder by database type.
func (d *Database) Select(table Tables, whereClause string) Result {
	result := Result{}

	return result
}

// Insert query builder by database type
func (d *Database) Insert(table Tables, data interface{}) Result {
	result := Result{}

	return result
}

// Update query builder by database type
func (d *Database) Update(table Tables, whereClause string, data interface{}) Result {
	result := Result{}

	return result
}

// Delete query build by database type
func (d *Database) Delete(table Tables, whereClause string) Result {
	result := Result{}

	return result
}
