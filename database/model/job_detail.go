// Copyright 2019 StreetByters Community
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
	"github.com/streetbyters/agente/database"
	"github.com/streetbyters/agente/model"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// JobDetail user defined background job detail structure
type JobDetail struct {
	database.DBInterface `json:"-"`
	ID                   int64    `db:"id" json:"id"`
	NodeID               int64    `db:"node_id" json:"node_id" foreign:"fk_ra_job_details_node_id" validate:"required"`
	JobID                int64    `db:"job_id" json:"job_id" foreign:"fk_ra_job_details_job_id" validate:"required"`
	SourceUserID         zero.Int `db:"source_user_id" foreign:"fk_ra_job_details_source_user_id" json:"source_user_id"`

	Code       string        `db:"code" json:"code" validate:"required,gte=3,lte=64"`
	Name       string        `db:"name" json:"name" validate:"required,gte=3,lte=200"`
	Type       model.JobType `db:"type" json:"type"`
	Detail     zero.String   `db:"detail" json:"detail"`
	Before     bool          `db:"before" json:"before"`
	BeforeJobs zero.String   `db:"before_jobs" json:"before_jobs"`
	After      bool          `db:"after" json:"after"`
	AfterJobs  zero.String   `db:"after_jobs" json:"after_jobs"`

	ScriptFile zero.String `db:"script_file" json:"script_file"`
	Script     zero.String `db:"script" json:"script"`

	InsertedAt time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewJobDetail generate user defined background job detail
func NewJobDetail() *JobDetail {
	return &JobDetail{}
}

// TableName job detail database table name
func (d JobDetail) TableName() string {
	return "ra_job_details"
}

// ToJSON job detail structure to json string
func (d JobDetail) ToJSON() string {
	return database.ToJSON(d)
}
