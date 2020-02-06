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

package errors

func New(text string, args ...interface{}) error {
	e := &PluggableError{s: text}

	if len(args) >= 1 {
		e.Status = args[0].(int)
	}

	if len(args) >= 2 {
		e.Detail = args[1].(string)
	}

	if len(args) >= 3 {
		e.Errors = args[2].(interface{})
	}

	return e
}

type PluggableError struct {
	Status int
	Errors interface{}
	Detail string
	s      string
}

func (e *PluggableError) Error() string {
	return e.s
}
