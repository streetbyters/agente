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

import (
	"errors"
)

// Scheduler application job scheduler
type Scheduler struct {
	App *App
	Package SchedulerInterface
}

// SchedulerInterface application job scheduler interface
type SchedulerInterface interface {
	Start()
	List()
	Add()
	Update(id int64)
	Delete(id int64)
	Run()
	Stop()
}

// NewScheduler building job scheduler
func NewScheduler(app *App) *Scheduler {
	s := &Scheduler{App: app}

	switch app.Config.Scheduler {
	case "JobRunner":
		s.Package = &SchedulerJobRunner{Scheduler: s}
		break
	default:
		panic(errors.New("unknown scheduler type"))
	}

	return s
}
