// Copyright 2019 Abdulkadir Dilsiz
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
	"github.com/akdilsiz/agente/database"
	"github.com/akdilsiz/agente/model"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// File job scripts database structure
type File struct {
	database.DBInterface `json:"-"`
	ID                   int64      `db:"id" json:"id"`
	NodeID               int64      `db:"node_id" json:"node_id" foreign:"fk_ra_files_node_id"`
	ParentID             zero.Int   `db:"parent_id" json:"parent_id" foreign:"fk_ra_files_parent_id"`
	JobID                zero.Int   `db:"job_id" json:"job_id" foreign:"fk_ra_files_job_id"`
	Dir                  string     `db:"dir" json:"dir" validate:"required"`
	File                 string     `db:"file" json:"file" validate:"required"`
	Type                 model.Node `db:"type" json:"type"`
	InsertedAt           time.Time  `db:"inserted_at" json:"inserted_at"`
	UpdatedAt            time.Time  `db:"updated_at" json:"updated_at"`
}

// NewFile generate file structure
func NewFile() *File {
	return &File{Type: model.Worker}
}

// TableName file structure database table name
func (d *File) TableName() string {
	return "ra_files"
}

// ToJSON file structure to json string
func (d *File) ToJSON() string {
	return database.ToJSON(d)
}
