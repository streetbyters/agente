//
// Copyright 2019 Abdulkadir DILSIZ <TransferChain>
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

package response

import (
	"encoding/json"
	"github.com/akdilsiz/release-agent/model"
)

type Error struct {
	model.Response			`json:"-"`
	Errors		interface{}		`json:"errors"`
	Detail		string			`json:"detail"`
}

func (r Error) ToJson() string {
	body, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(body)
}
