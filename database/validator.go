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
	"github.com/lib/pq"
	"github.com/streetbyters/agente/utils"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
)

var validate = validator.New()

// Tag error constraint structure
type Tag struct {
	Name       string
	Constraint string
}

// ValidateStruct struct validator
func ValidateStruct(r interface{}) (map[string]string, error) {
	err := validate.Struct(r)
	errors := map[string]string{}

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Param() != "" {
				errors[utils.ToSnakeCase(err.Field())] = err.ActualTag() + ": " + err.Param()
			} else {
				errors[utils.ToSnakeCase(err.Field())] = err.ActualTag()
			}

		}
	}

	return errors, err
}

// ValidateConstraint sql violation error validator
func ValidateConstraint(err error, r interface{}) (map[string]string, error) {
	errs := map[string]string{}
	if pgerr, ok := err.(*pq.Error); ok {
		t := reflect.TypeOf(r).Elem()

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i).Tag

			var tags []Tag

			if tag, ok := field.Lookup("unique"); ok && tag == pgerr.Constraint {
				uTags := strings.Split(tag, ",")
				for _, _t := range uTags {
					if _t == pgerr.Constraint {
						tags = append(tags, struct {
							Name       string
							Constraint string
						}{Name: "unique", Constraint: pgerr.Constraint})
					}
				}
			}
			errs = constraintErrors(t.Field(i), pgerr, tags, errs)

			if tag, ok := field.Lookup("foreign"); ok && tag == pgerr.Constraint {
				tags = append(tags, struct {
					Name       string
					Constraint string
				}{Name: "foreign", Constraint: pgerr.Constraint})
			}

			errs = constraintErrors(t.Field(i), pgerr, tags, errs)
		}
	}

	return errs, err
}

func constraintErrors(field reflect.StructField, dbError *pq.Error, tags []Tag, errs map[string]string) map[string]string {
	for range tags {
		var msg string

		switch string(dbError.Code) {
		case string(UniqueViolation):
			msg = string(UniqueError)
			break
		case string(ForeignKeyViolation):
			msg = string(NotExistsError)
			break
		}

		errs[utils.ToSnakeCase(field.Name)] = msg
	}

	return errs
}
