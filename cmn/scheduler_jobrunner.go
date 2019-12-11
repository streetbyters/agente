//
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

import "github.com/bamzi/jobrunner"

// SchedulerJobRunner jobrunner package adapter
type SchedulerJobRunner struct {
	SchedulerInterface `json:"-"`
	*Scheduler
}

// Start jobrunner
func (s SchedulerJobRunner) Start() {
	jobrunner.Start()
}

// List jobrunner jobs
func (s SchedulerJobRunner) List() {

}

// Add jobrunner job
func (s SchedulerJobRunner) Add() {

}

// Update jobrunner job
func (s SchedulerJobRunner) Update(id int64) {

}

// Delete jobrunner job
func (s SchedulerJobRunner) Delete(id int64) {

}

// Run jobrunner job
func (s SchedulerJobRunner) Run() {

}

// Stop jobrunner
func (s SchedulerJobRunner) Stop() {
	jobrunner.Stop()
}
