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

// DB enum type
type DB string

// Postgres Enum
const Postgres DB = "postgres"

// Unknown enum. This enum just test mode
const Unknown DB = "unknown"

// MODE type for application
type MODE string

// Dev Development mode enum
const Dev MODE = "dev"

// Test mode enum
const Test MODE = "test"

// Prod Production model enum
const Prod MODE = "prod"

// JobType received message type
type JobType string

const (
	// NewRelease job type
	NewRelease JobType = "new_release"
	// Start job type
	Start JobType = "start"
	// Restart job type
	Restart JobType = "restart"
	// Shutdown job type
	Shutdown JobType = "shutdown"
	// Other job type
	Other JobType = "other"
)

// Node type for worker
type Node string

const (
	// Worker node type
	Worker Node = "worker"
	// Master node type
	Master Node = "master"
)

// Process type for file operation
type Process string

const (
	// Insert process for file operation
	Insert Process = "insert"
	// Update process for file operation
	Update Process = "update"
	// Distributing process for file operation
	Distributing Process = "distributing"
)
