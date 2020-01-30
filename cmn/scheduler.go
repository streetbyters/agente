// Copyright 2019 Forgolang Community
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
	"errors"
	"github.com/forgolang/agente/utils"
)

// Scheduler application job scheduler
type Scheduler struct {
	App     *App
	Package SchedulerInterface
}

// SchedulerJob General scheduler job struct
type SchedulerJob struct {
	Data interface{} `json:"data"`
}

// Packages All defined scheduler packages
func (s *Scheduler) Packages() ([]string, map[string]SchedulerInterface) {
	packages := make(map[string]SchedulerInterface)

	packages["GoCron"] = &SchedulerGoCron{Scheduler: s}

	return []string{"GoCron"}, packages
}

// SchedulerInterface application job scheduler interface
type SchedulerInterface interface {
	Up()
	Down()
	Start()
	List()
	Add(args ...interface{})
	Update(args ...interface{})
	Delete(args ...interface{})
	Run()
	Stop()
}

// NewScheduler building job scheduler
func NewScheduler(app *App) *Scheduler {
	s := &Scheduler{App: app}

	names, packages := s.Packages()

	if ok, _ := utils.InArray(app.Config.Scheduler, names); !ok {
		panic(errors.New("undefined scheduler"))
	}

	s.Package = packages[app.Config.Scheduler]
	s.Package.Up()
	return s
}
