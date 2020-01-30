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

package utils

import "gopkg.in/go-playground/validator.v9"

var validate = validator.New()

// ValidateStruct struct validator
func ValidateStruct(r interface{}) (map[string]string, error) {
	err := validate.Struct(r)
	errors := map[string]string{}

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Param() != "" {
				errors[ToSnakeCase(err.Field())] = err.ActualTag() + ": " + err.Param()
			} else {
				errors[ToSnakeCase(err.Field())] = err.ActualTag()
			}

		}
	}

	return errors, err
}
