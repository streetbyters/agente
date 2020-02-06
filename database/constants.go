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

package database

// Postgres type for postgres contants
type Postgres string

// Violation type for violation constants
type Violation string

const (
	// UniqueViolation sql unique violation code
	UniqueViolation Postgres = "23505"
	// ForeignKeyViolation sql foreign key violation code
	ForeignKeyViolation Postgres = "23503"
	// NotNullViolation sql not null violation code
	NotNullViolation Postgres = "23502"
)

const (
	// UniqueError unique violation error string
	UniqueError Violation = "has been already taken"
	// NotExistsError not exists violation error string
	NotExistsError Violation = "does not exists"
	// NotNullError not null violation error string
	NotNullError Violation = "is not null"
)
