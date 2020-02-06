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

package cmn

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
	model2 "github.com/streetbyters/agente/database/model"
)

// SchedulerGoCron gocron package adapter
type SchedulerGoCron struct {
	SchedulerInterface `json:"-"`
	*Scheduler
	GoCron *gocron.Scheduler
	Jobs   []Job
}

// Up gocron scheduler
func (s *SchedulerGoCron) Up() {
	s.GoCron = gocron.NewScheduler()
	jobDetail := model2.NewJobDetail()
	job := model2.NewJob()
	var jobs []model2.Job
	res := s.Scheduler.App.Database.QueryWithModel(fmt.Sprintf("SELECT * FROM %s"+
		" AS j", job.TableName()),
		&jobs)
	var jobIDs []int64
	for _, j := range jobs {
		jobIDs = append(jobIDs, j.ID)
	}

	var rJobs []model2.Job
	jobDetails := make([]model2.JobDetail, 0)
	details := make(map[int64]model2.JobDetail)
	if res.Error == nil {
		query, args, _ := sqlx.In(fmt.Sprintf("SELECT d.* FROM %s AS d"+
			" LEFT OUTER JOIN %s AS d2 ON d.job_id = d2.job_id AND d.id < d2.id"+
			" WHERE d2.id IS NULL AND d.job_id IN (?)", jobDetail.TableName(), jobDetail.TableName()),
			jobIDs)
		query = s.Scheduler.App.Database.DB.Rebind(query)
		s.Scheduler.App.Database.QueryWithModel(query, &jobDetails, args...)

		for _, d := range jobDetails {
			details[d.JobID] = d
		}

		for _, j := range jobs {
			if val, ok := details[j.ID]; ok {
				j.Detail = &val
			}
			rJobs = append(rJobs, j)
		}
	}
}

// Start gocron
func (s *SchedulerGoCron) Start() {
	s.GoCron.Start()
}

// List gcron jobs
func (s *SchedulerGoCron) List() {

}

// Add gcron job
func (s *SchedulerGoCron) Add(args ...interface{}) {

}

// Update gocron job
func (s *SchedulerGoCron) Update(args ...interface{}) {

}

// Delete gocron job
func (s *SchedulerGoCron) Delete(args ...interface{}) {

}

// Run gocron job
func (s *SchedulerGoCron) Run() {

}

// Stop gocron job
func (s *SchedulerGoCron) Stop() {

}

// Down gocron kill
func (s *SchedulerGoCron) Down() {
	s.GoCron.Clear()
	s.GoCron = nil
}
