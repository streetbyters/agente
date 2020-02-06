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

package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fate-lovely/phi"
	"github.com/streetbyters/agente/model"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

// JWTAuth authentication mechanism
type JWTAuth struct {
	API    *API
	Method *jwt.SigningMethodHMAC
	Secret string
	Expire int64
}

// NewJWTAuth generate jwt auth
func NewJWTAuth(api *API) *JWTAuth {
	return &JWTAuth{
		API:    api,
		Method: jwt.SigningMethodHS256,
		Secret: api.App.Config.SecretKey,
	}
}

// Generate generate jwt token with mapClaims
func (a JWTAuth) Generate(args ...interface{}) (string, error) {
	a.Expire = time.Now().UTC().Unix() + int64(time.Second*25)

	claims := jwt.MapClaims{
		"id":  args[0].(int64),
		"exp": a.Expire,
	}

	token := jwt.NewWithClaims(a.Method, claims)

	tokenString, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Parse token string parse mapClaims expire check
func (a JWTAuth) Parse(tokenString string) (map[string]interface{}, int) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.Secret), nil
	})

	if err != nil {
		return map[string]interface{}{}, -1
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if time.Now().UTC().Unix() > int64(claims["exp"].(float64)) {
			return map[string]interface{}{}, 0
		}

		return claims, 1
	}

	return map[string]interface{}{}, -2
}

// Verify verify bearer token in requests
func (a JWTAuth) Verify(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		h := ctx.Request.Header.Peek("authorization")

		if len(string(h)) < 7 && strings.ToUpper(string(h)[0:6]) != "BEARER" {
			a.API.JSONResponse(ctx, model.ResponseError{
				Detail: "the request sent did not consist a valid token entry.",
			}, fasthttp.StatusForbidden)
			return
		}

		claims, err := a.Parse(string(h)[7:])

		switch err {
		case 0:
			a.API.JSONResponse(ctx, model.ResponseError{
				Detail: "token expire",
			}, fasthttp.StatusUnauthorized)
			return
		case -1, -2:
			a.API.JSONResponse(ctx, model.ResponseError{
				Detail: "the token supplied could not be validated.",
			}, fasthttp.StatusForbidden)
			return
		default:
			a.API.Auth.ID = int64(claims["id"].(float64))

			next(ctx)
			break
		}
	}
}
