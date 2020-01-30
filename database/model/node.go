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
	"github.com/forgolang/agente/database"
	"github.com/forgolang/agente/model"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// Node application type structure
type Node struct {
	database.DBInterface `json:"-"`
	ID                   int64       `db:"id" json:"id"`
	Name                 string      `db:"name" json:"name" validate:"required,gte=3,lte=200"`
	Code                 string      `db:"code" json:"code" unique:"ra_nodes_code_unique_index" validate:"required,gte=3,lte=200"`
	Detail               zero.String `db:"detail" json:"detail"`
	Type                 model.Node  `db:"type" json:"type"`
	InsertedAt           time.Time   `db:"inserted_at" json:"inserted_at"`
	UpdatedAt            time.Time   `db:"updated_at" json:"updated_at"`
}

// NewNode generate node structure
func NewNode() *Node {
	return &Node{Type: "worker"}
}

// TableName node database table name
func (d *Node) TableName() string {
	return "ra_nodes"
}

// ToJSON node structure to json string
func (d *Node) ToJSON() string {
	return database.ToJSON(d)
}
