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

import "github.com/jasonlvhit/gocron"

// SchedulerGoCron gocron package adapter
type SchedulerGoCron struct {
	SchedulerInterface `json:"-"`
	*Scheduler
	GoCron *gocron.Scheduler
}

// Up gocron scheduler
func (s *SchedulerGoCron) Up() {
	s.GoCron = gocron.NewScheduler()
}

// Start gocron
func (s *SchedulerGoCron) Start() {

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
