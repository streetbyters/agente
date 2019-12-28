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

package model

import (
	"github.com/akdilsiz/agente/database"
	"time"
)

// UserPassphraseInvalidation authentication access token invalidation struct
type UserPassphraseInvalidation struct {
	database.DBInterface `json:"-"`
	PassphraseID         int64     `db:"passphrase_id" json:"passphrase_id" foreign:"fk_ra_user_passphrases_passphrase_id" unique:"ra_user_passphrase_invalidations_pkey"`
	NodeID               int64     `db:"node_id" json:"node_id" foreign:"fk_ra_user_passphrases_invalidations_node_id" validate:"required"`
	SourceUserID         int64     `db:"source_user_id" json:"source_user_id" foreign:"fk_ra_user_passphrases_source_user_id"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewUserPassphraseInvalidation generate authentication access token invalidation
func NewUserPassphraseInvalidation() *UserPassphraseInvalidation {
	return &UserPassphraseInvalidation{}
}

// TableName user_passphrase_invalidation database table name
func (d UserPassphraseInvalidation) TableName() string {
	return "ra_user_passphrase_invalidations"
}

func (d UserPassphraseInvalidation) ToJSON() string {
	return database.ToJSON(d)
}
