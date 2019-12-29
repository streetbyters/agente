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

package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/akdilsiz/agente/database"
	"github.com/akdilsiz/agente/model"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// FileLog log for file database structure
type FileLog struct {
	database.DBInterface `json:"-"`
	ID                   int64         `db:"id" json:"id"`
	NodeID               int64         `db:"node_id" json:"node_id" foreign:"fk_ra_file_logs_node_id"`
	FileID               int64         `db:"file_id" json:"file_id" foreign:"fk_ra_file_logs_file_id"`
	SourceUserID         zero.Int      `db:"source_user_id" json:"source_user_id" foreign:"fk_ra_file_logs_source_user_id"`
	Type                 model.Process `db:"type" json:"type"`
	Data                 FileLogData   `db:"data" json:"data,omitempty"`
	InsertedAt           time.Time     `db:"inserted_at" json:"inserted_at"`
}

// NewFileLog generate file log structure
func NewFileLog(fileID int64) *FileLog {
	return &FileLog{FileID: fileID}
}

// TableName upload log structure database table name
func (d *FileLog) TableName() string {
	return "ra_file_logs"
}

// ToJSON file log structure to json string
func (d *FileLog) ToJSON() string {
	return database.ToJSON(d)
}

// FileLogData jsonb structure
type FileLogData struct {
	Dir  string     `json:"dir,omitempty"`
	File string     `json:"file,omitempty"`
	Type model.Node `json:"type,omitempty"`
}

// Value file log data driver.Valuer
func (a FileLogData) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan file log data sql.Scanner
func (a *FileLogData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
