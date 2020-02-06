// Copyright 2019 Street Byters Community
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
	"github.com/forgolang/agente/database"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// Job user defined background job structure
type Job struct {
	database.DBInterface `json:"-"`
	ID                   int64      `db:"id" json:"id"`
	NodeID               int64      `db:"node_id" json:"node_id" foreign:"fk_ra_jobs_node_id" validate:"required"`
	SourceUserID         zero.Int   `db:"source_user_id" json:"source_user_id" foreign:"fk_ra_jobs_source_user_id"`
	InsertedAt           time.Time  `db:"inserted_at" json:"inserted_at"`
	Detail               *JobDetail `json:"detail"`
}

// NewJob generate user defined background job
func NewJob() *Job {
	return &Job{}
}

// TableName job database table name
func (d *Job) TableName() string {
	return "ra_jobs"
}

// ToJSON job structure to json string
func (d *Job) ToJSON() string {
	return database.ToJSON(d)
}
