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

import (
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v3/zero"
	"reflect"
	"time"
)

// Change database model compare and change
type Change struct {
	Key     string
	Name    string
	String  string
	Int     int64
	Bool    bool
	Decimal decimal.Decimal
}

// GetChanges Get repo model changes
func GetChanges(m interface{}, c interface{}, typs ...string) ([]Change, []string, map[string]interface{}) {
	var changes []Change
	var keys []string
	namedParams := make(map[string]interface{})

	if len(typs) <= 0 {
		return changes, keys, namedParams
	}
	var typ string
	typ = typs[0]

	t := reflect.TypeOf(m).Elem()
	t2 := reflect.TypeOf(c).Elem()

	values := reflect.ValueOf(m).Elem()
	values2 := reflect.ValueOf(c).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name != "DBInterface" {
			changes, keys, namedParams = getChanges(typ, t2, field, values, values2, changes, keys, namedParams)
		}
	}

	for k, v := range namedParams {
		if k == "" || v == "" {
			delete(namedParams, k)
		}
	}

	return changes, keys, namedParams
}

func getChanges(args ...interface{}) ([]Change, []string, map[string]interface{}) {
	typ := args[0].(string)
	t2 := args[1].(reflect.Type)
	field := args[2].(reflect.StructField)
	values := args[3].(reflect.Value)
	values2 := args[4].(reflect.Value)
	changes := args[5].([]Change)
	keys := args[6].([]string)
	namedParams := args[7].(map[string]interface{})

	if field2, ok := t2.FieldByName(field.Name); ok && field2.Tag.Get("db") != "" && field.Name != "ID" {
		var change Change

		val := values.FieldByName(field.Name)
		val2 := values2.FieldByName(field2.Name)
		switch field.Type.Kind() {
		case reflect.Struct:
			switch field.Type {
			case reflect.TypeOf(zero.String{}):
				if val.FieldByName("String").String() != val2.FieldByName("String").String() &&
					val2.FieldByName("String").String() != "" {
					change.Key = field.Name
					change.Name = field.Tag.Get("db")
					change.String = val2.FieldByName("String").String()
					changes = append(changes, change)
					keys = append(keys, change.Name)
					namedParams[change.Name] = change.String
					val.Set(val2)
				}
				break
			case reflect.TypeOf(zero.Int{}):
				if val.FieldByName("Int64").Int() != val2.FieldByName("Int64").Int() &&
					val2.FieldByName("Int64").Int() != 0 {
					change.Key = field.Name
					change.Name = field.Tag.Get("db")
					change.Int = val2.FieldByName("Int64").Int()
					changes = append(changes, change)
					keys = append(keys, change.Name)
					namedParams[change.Name] = change.Int
					val.FieldByName("Int64").SetInt(val2.FieldByName("Int64").Int())
				} else if typ == "insert" {
					namedParams[change.Name] = val2.FieldByName("Int64").Int()
					val.Set(val2)
				}
				break
			case reflect.TypeOf(zero.Bool{}):
				change.Key = field.Name
				change.Name = field.Tag.Get("db")
				change.Bool = val2.FieldByName("Bool").Bool()
				changes = append(changes, change)
				keys = append(keys, change.Name)
				namedParams[change.Name] = change.Bool
				val.Set(val)
				break
			case reflect.TypeOf(decimal.NullDecimal{}):
				change.Key = field.Name
				change.Name = field.Tag.Get("db")
				change.Decimal = val2.Field(0).Interface().(decimal.Decimal)
				changes = append(changes, change)
				keys = append(keys, change.Name)
				namedParams[change.Name] = change.Decimal
				val.Set(val2)
				break
			case reflect.TypeOf(zero.Time{}):
				if val.Interface() != val2.Interface() {
					change.Key = field.Name
					change.Name = field.Tag.Get("db")
					change.String = val2.Interface().(zero.Time).Time.String()
					changes = append(changes, change)
					keys = append(keys, change.Name)
					namedParams[change.Name] = change.String
					val.Set(val2)
				}
				break
			case reflect.TypeOf(time.Time{}):
				if val.Interface() != val2.Interface() {
					change.Key = field.Name
					change.Name = field.Tag.Get("db")
					change.String = val2.Interface().(time.Time).String()
					changes = append(changes, change)
					keys = append(keys, change.Name)
					namedParams[change.Name] = change.String
					val.Set(val2)
				}
				break
			}
		case reflect.String:
			if val.String() != val2.String() {
				change.Key = field.Name
				change.Name = field.Tag.Get("db")
				change.String = val2.String()
				changes = append(changes, change)
				keys = append(keys, change.Name)
				namedParams[change.Name] = change.String
				val.Set(val2)
			}
			break
		case reflect.Int64:
			if val.Int() != val2.Int() {
				change.Key = field.Name
				change.Name = field.Tag.Get("db")
				change.Int = val2.Int()
				changes = append(changes, change)
				keys = append(keys, change.Name)
				namedParams[change.Name] = change.Int
				val.Set(val2)
			}
			break
		case reflect.Int:
			if val.Int() != val2.Int() {
				change.Key = field.Name
				change.Name = field.Tag.Get("db")
				change.Int = val2.Int()
				changes = append(changes, change)
				keys = append(keys, change.Name)
				namedParams[change.Name] = change.Int
				val.Set(val2)
			}
			break
		case reflect.Bool:
			change.Key = field.Name
			change.Name = field.Tag.Get("db")
			change.Bool = val2.Bool()
			changes = append(changes, change)
			keys = append(keys, change.Name)
			namedParams[change.Name] = change.Bool
			val.Set(val2)
			break
		}
	}

	return changes, keys, namedParams
}
