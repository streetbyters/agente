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

package database

import (
	"encoding/json"
	"errors"
	"fmt"
	pluggableError "github.com/akdilsiz/agente/errors"
	"github.com/akdilsiz/agente/model"
	"github.com/akdilsiz/agente/utils"
	_ "github.com/go-sql-driver/mysql" // Mysql Driver
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // Postgres Driver
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sort"
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
	Config    *model.Config
	Type      model.DB
	DB        *sqlx.DB
	Tx        *sqlx.Tx
	Logger    *utils.Logger
	Error     error
	Reset     bool
	QueryType string
}

// Force raise panic database query result is nil
func (d Database) Force() Database {
	if d.Error != nil {
		switch d.QueryType {
		case "row":
			panic(pluggableError.New("not found", fasthttp.StatusNotFound))
		case "insert", "update":
			panic(pluggableError.New("incorrect given parameters", fasthttp.StatusUnprocessableEntity))
		}
	}

	defer func() {
		d.QueryType = ""
	}()

	return d
}

// DBInterface database model interface
type DBInterface interface {
	TableName() string
	ToJSON() string
}

// ToJSON Converting models belonging to DBInterface to json string
func ToJSON(model DBInterface) string {
	b, err := json.Marshal(model)
	if err != nil {
		return ""
	}
	return string(b)
}

// Tx transaction for database queries
type Tx struct {
	DB *Database
}

// Result structure for database query results
type Result struct {
	QueryType string
	Rows      []interface{}
	Count     int64
	Error     error
}

// TODO: parameters check
// Force raise panic database query result is nil
func (r Result) Force() Result {
	if r.Error != nil {
		switch r.QueryType {
		case "row":
			panic(pluggableError.New("not found", fasthttp.StatusNotFound))
		}
	}

	defer func() {
		r.QueryType = ""
	}()

	return r
}

// NewDB building database
func NewDB(config *model.Config) (*Database, error) {
	database := &Database{}
	database.Config = config

	switch config.DB {
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
	default:
		return nil, errors.New("unsupported database")
	}

	return database, nil
}

// DropDB Drop database schemas
func DropDB(database *Database) error {
	var err error
	switch database.Type {
	case model.Postgres:
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

	database.reset()
	switch database.Type {
	case model.Postgres:
		err = migrationUp(database)
		break
	}

	return err
}

func (d *Database) reset() {
	switch d.Type {
	case model.Postgres:
		if d.Reset {
			//d.DB.Exec("DROP SCHEMA public CASCADE;")
			//d.DB.Exec("CREATE SCHEMA public;")
			//d.DB.Exec("GRANT ALL ON SCHEMA public TO postgres;")
			//d.DB.Exec("GRANT ALL ON SCHEMA public TO public;")
			//d.DB.Exec("")
			files := migrationFiles(d, "down")
			for _, f := range files {
				if _, err := d.DB.Exec(f.Data); err != nil {
					panic(err)
				}
			}
		}
		break
	}
}

type sqlS struct {
	Number int
	Name   string
	Data   string
}

func migrationFiles(db *Database, typ string) []sqlS {
	var sqls []sqlS

	var files []string

	files, _ = filepath.Glob(filepath.Join(db.Config.DBPath, "sql", string(db.Config.DB), "[0-9]*.[a-zA-Z_]*."+typ+".sql"))

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

	if typ == "down" {
		sort.Slice(sqls, func(i, j int) bool {
			return sqls[i].Number > sqls[j].Number
		})
	}

	return sqls
}

func migrationUp(db *Database) error {
	if err := baseMigrations(db); err != nil {
		return err
	}

	return newMigrations(db)
}

func baseMigrations(db *Database) error {
	var err error
	files := migrationFiles(db, "up")

	_, err = db.DB.Queryx("SELECT * FROM " + string(tMigration) + " AS m ORDER BY id ASC")
	if err != nil {
		if int(dbError(db, err)) == int(TableNotFound) {
			err = nil
			tx, _ := db.DB.Beginx()
			for _, f := range files {
				switch f.Name {
				case "01.postgres.up.sql":
					_, err = tx.Exec(f.Data)
					break
				}
			}

			if err != nil {
				tx.Rollback()
				return err
			}
			tx.Commit()
		}

		return nil
	}

	return err
}

func newMigrations(db *Database) error {
	var err error
	result := Result{}
	result = db.Query("SELECT * FROM " + string(tMigration) + " AS m ORDER BY id ASC")
	var lastMigration []interface{}
	if len(result.Rows) > 0 {
		lastMigration = result.Rows[:len(result.Rows)]
	}

	tx, err := db.DB.Beginx()
	files := migrationFiles(db, "up")

	for _, f := range files {
		switch f.Name {
		case "01.postgres.up.sql":
			break
		default:
			if len(lastMigration) > 0 {
				if f.Number > int(lastMigration[1].(int64)) {
					_, err = tx.Exec(f.Data)
					if err != nil {
						tx.Rollback()
						break
					}

					_, err = tx.Exec("INSERT INTO "+string(tMigration)+" ("+
						"number, name) VALUES ($1, $2)", f.Number, f.Name)
					if err == nil {
						db.Logger.LogInfo("Migrate: " + f.Name)
					}
				}
			} else {
				_, err = tx.Exec(f.Data)
				if err != nil {
					tx.Rollback()
					break
				}

				_, err = tx.Exec("INSERT INTO "+string(tMigration)+" ("+
					"number, name) VALUES ($1, $2)", f.Number, f.Name)
				if err == nil {
					db.Logger.LogInfo("Migrate: " + f.Name)
				}
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

// QueryWithModel database query builder with given model
func (d *Database) QueryWithModel(query string, target interface{}, params ...interface{}) Result {
	return d.query(query, target, params...)
}

// Query database query builder
func (d *Database) Query(query string, params ...interface{}) Result {
	return d.query(query, nil, params...)
}

func (d *Database) query(query string, target interface{}, params ...interface{}) Result {
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
	defer rows.Close()

	var arr reflect.Value
	var v reflect.Value
	if target != nil {
		arr = reflect.ValueOf(target).Elem()
		v = reflect.New(reflect.TypeOf(target).Elem().Elem())
	}

	for rows.Next() {
		if target != nil {
			if err := rows.StructScan(v.Interface()); err != nil {
				result.Error = err
				break
			}
			arr.Set(reflect.Append(arr, v.Elem()))
		} else {
			row, err := rows.SliceScan()
			if err != nil {
				result.Error = err
				break
			}
			result.Rows = append(result.Rows, row)
		}
	}
	if err := rows.Err(); err != nil {
		result.Error = err
		return result
	}

	return result
}

// QueryRowWithModel database row query builder with target model
func (d *Database) QueryRowWithModel(query string, target interface{}, params ...interface{}) Result {
	return d.queryRow(query, target, params...)
}

// QueryRow database row query builder
func (d *Database) QueryRow(query string, params ...interface{}) Result {
	return d.queryRow(query, nil, params...)
}

func (d *Database) queryRow(query string, target interface{}, params ...interface{}) Result {
	result := Result{}
	result.QueryType = "row"

	var err error
	var r interface{}
	var row *sqlx.Row

	if d.Tx != nil {
		row = d.Tx.QueryRowx(query, params...)
	} else {
		row = d.DB.QueryRowx(query, params...)
	}

	if target != nil {
		err = row.StructScan(target)
	} else {
		r, err = row.SliceScan()
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
func (t *Tx) Select(table string, whereClause string) Result {
	result := Result{}

	result.QueryType = "row"

	return result
}

// Insert query builder by database type
func (d *Database) Insert(m DBInterface, data interface{}, keys ...string) error {
	_, c1, _ := GetChanges(m, data, "insert")

	d.QueryType = "insert"

	str, _ := insertSQL(c1, m.TableName(), strings.Join(keys, ", "))

	var stmt *sqlx.NamedStmt
	var err error

	if d.Tx != nil {
		stmt, err = d.Tx.PrepareNamed(str)
	} else {
		stmt, err = d.DB.PrepareNamed(str)
	}

	if err != nil {
		d.Error = err
		return err
	}

	if err := stmt.QueryRowx(data).StructScan(data); err != nil {
		d.Error = err
		if d.Tx != nil {
			d.rollback()
		}
		return err
	}

	if err := stmt.Close(); err != nil {
		d.Error = err
		if d.Tx != nil {
			d.rollback()
		}

		return err
	}

	return nil
}

// Update query builder by database type
func (d *Database) Update(table string, whereClause string, data interface{}) Result {
	result := Result{}

	return result
}

// Delete query build by database type
func (d *Database) Delete(table string, whereClause string, args ...interface{}) Result {
	result := Result{}
	d.QueryType = "row"

	if d.Tx != nil {
		res, err := d.Tx.Exec(fmt.Sprintf("DELETE FROM %s", table) + " WHERE " + whereClause, args...)
		result.Error = err
		if err != nil {
			return result
		}

		if i, _ := res.RowsAffected(); i <= 0 {
			result.Error = errors.New("do not found row affected")
		}

		return result
	}

	res, err := d.DB.Exec(fmt.Sprintf("DELETE FROM %s", table) + " WHERE " + whereClause, args...)
	result.Error = err
	if err != nil {
		return result
	}
	if i, _ := res.RowsAffected(); i <= 0 {
		result.Error = errors.New("do not found row affected")
	}

	return result
}

func insertSQL(columns []string, tableName string, keyColumn string, args ...interface{}) (string, error) {
	tmplStr := `insert into {{.TableName}} (` +
		`{{$putComa := false}}` +
		`{{- range $i, $f := .Columns}}` +
		`{{if $putComa}}, {{end}}{{$f}}{{$putComa = true}} ` +
		`{{- end}}` +
		`) values (` +
		`{{$putComa := false}}` +
		`{{- range $i, $f := .Columns}}` +
		`{{if $putComa}}, {{end}}:{{$f}}{{$putComa = true}} ` +
		`{{- end}}` +
		`) ` +
		`{{if ne .KeyColumn ""}}returning {{.KeyColumn}}{{end}}`

	data := struct {
		TableName string
		Columns   []string
		KeyColumn string
	}{
		TableName: tableName,
		Columns:   columns,
		KeyColumn: keyColumn,
	}

	return utils.ParseAndExecTemplateFromString(tmplStr, data)
}

func updateSQL(columns []string, tableName string, whereClause string, keyColumn string) (string, error) {
	tmplStr := `update {{.TableName}} set ` +
		`{{$putComa := false}} ` +
		`{{- range $i, $f := .Columns}}` +
		`{{if $putComa}}, {{end}}{{$f}} = :{{$f}}{{$putComa = true}} ` +
		`{{- end}} ` +
		`where {{.WhereClause}} ` +
		`{{if ne .KeyColumn ""}}returning {{.KeyColumn}}{{end}}`

	data := struct {
		TableName   string
		Columns     []string
		KeyColumn   string
		WhereClause string
	}{
		TableName:   tableName,
		Columns:     columns,
		KeyColumn:   keyColumn,
		WhereClause: whereClause,
	}

	return utils.ParseAndExecTemplateFromString(tmplStr, data)
}
