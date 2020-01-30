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

package model

import (
	"errors"
	"github.com/forgolang/agente/utils"
	"strings"
)

// Pagination request pagination field parser
type Pagination struct {
	Model
	Limit      int `validate:"lte=40"`
	Offset     int64
	OrderBy    string `validate:"oneof=asc desc"`
	OrderField string
}

// NewPagination generate pagination struct with default values
func NewPagination() Pagination {
	return Pagination{
		Limit:      40,
		Offset:     0,
		OrderBy:    "desc",
		OrderField: "id",
	}
}

// Validate pagination fields validator
func (m Pagination) Validate(orderFields ...string) (map[string]string, error) {
	exists, _ := utils.InArray(m.OrderField, orderFields)
	if len(orderFields) > 0 && !exists {
		return map[string]string{"order_field": "is not valid"},
			errors.New("the value entered should be " +
				strings.Join(orderFields, ", ") + " one of values.")
	}
	return utils.ValidateStruct(m)
}
