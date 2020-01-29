package database

import (
	"github.com/akdilsiz/agente/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3/zero"
	"testing"
	"time"
)

type TestStruct struct {
	Name   string `validate:"required"`
	Code   string
	Detail string `validate:"gte=5"`
}

// Node application type structure
type TestNode struct {
	DBInterface `json:"-"`
	ID          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name" validate:"required,gte=3,lte=200"`
	Code        string      `db:"code" json:"code" unique:"ra_nodes_code_unique_index" validate:"required,gte=3,lte=200"`
	Detail      zero.String `db:"detail" json:"detail"`
	Type        model.Node  `db:"type" json:"type"`
	InsertedAt  time.Time   `db:"inserted_at" json:"inserted_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updated_at"`
}

func (m *TestNode) TableName() string {
	return "ra_nodes"
}

type TestJobDetail struct {
	DBInterface  `json:"-"`
	ID           int64    `db:"id" json:"id"`
	NodeID       int64    `db:"node_id" json:"node_id"`
	JobID        int64    `db:"job_id" json:"job_id" foreign:"fk_ra_job_details_job_id" validate:"required"`
	SourceUserID zero.Int `db:"source_user_id" foreign:"fk_ra_job_details_source_user_id" json:"source_user_id"`

	Code       string        `db:"code" json:"code" validate:"required,gte=3,lte=64"`
	Name       string        `db:"name" json:"name" validate:"required,gte=3,lte=200"`
	Type       model.JobType `db:"type" json:"type"`
	Detail     zero.String   `db:"detail" json:"detail"`
	Before     bool          `db:"before" json:"before"`
	BeforeJobs zero.String   `db:"before_jobs" json:"before_jobs"`
	After      bool          `db:"after" json:"after"`
	AfterJobs  zero.String   `db:"after_jobs" json:"after_jobs"`

	ScriptFile zero.String `db:"script_file" json:"script_file"`
	Script     zero.String `db:"script" json:"script"`

	InsertedAt time.Time `db:"inserted_at" json:"inserted_at"`
}

func (m *TestJobDetail) TableName() string {
	return "ra_job_details"
}

type TestUser struct {
	DBInterface    `json:"-"`
	ID             int64     `db:"id" json:"id"`
	NodeID         int64     `db:"node_id" json:"node_id"`
	Username       string    `db:"username" json:"username" unique:"ra_users_username_unique_index" validate:"required"`
	PasswordDigest string    `db:"password_digest" json:"-"`
	Password       string    `db:"-" json:"password" validate:"required"`
	Email          string    `db:"email" json:"email" unique:"ra_users_email_unique_index" validate:"required,email"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	InsertedAt     zero.Time `db:"inserted_at" json:"inserted_at"`
	UpdatedAt      zero.Time `db:"updated_at" json:"updated_at"`
}

func (d TestUser) TableName() string {
	return "ra_users"
}

func TestValidateStruct(t *testing.T) {
	testStruct := new(TestStruct)
	testStruct.Code = "code"
	testStruct.Detail = "det"

	errs, err := ValidateStruct(testStruct)
	assert.NotNil(t, err)
	assert.Equal(t, errs["detail"], "gte: 5")
	assert.Equal(t, errs["name"], "required")
}

func TestValidateConstraint(t *testing.T) {
	appPath = dirs[0]

	// Open postgres db connection
	config := &model.Config{
		DBPath: appPath,
		Mode:   model.Test,
		DB:     model.Postgres,
		DBName: "agente_test",
		DBHost: "127.0.0.1",
		DBPort: 5432,
		DBUser: "agente",
		DBPass: "123456",
		DBSsl:  "disable",
	}

	db, err := NewDB(config)
	if err != nil {
		t.Fatal(err)
	}
	db.Logger = logger
	db.Reset = true
	DropDB(db)
	InstallDB(db)

	node := new(TestNode)
	node.Name = "node1"
	node.Code = "node1"
	err = db.Insert(new(TestNode), node, "id", "inserted_at")
	assert.Nil(t, err)

	// ForeignKey constraint error
	detail := new(TestJobDetail)
	detail.NodeID = node.ID
	detail.JobID = int64(999999999)
	detail.Code = "job3"
	detail.Name = "Test Job"
	detail.Detail.SetValid("Test Job Detail")
	detail.Before = true
	detail.BeforeJobs.SetValid("job4,job5")
	detail.After = true
	detail.AfterJobs.SetValid("job1,job2")
	detail.Script.SetValid(`
		#!/bin/sh
		echo "Test Job running"
	`)

	err = db.Insert(new(TestJobDetail), detail)
	assert.NotNil(t, err)
	errs, err := ValidateConstraint(err, detail)
	assert.Equal(t, errs["job_id"], "does not exists")

	// Unique constraint erro
	user := new(TestUser)
	user.NodeID = node.ID
	user.Username = "akdilsiz2"
	user.PasswordDigest = "asdasdasd"
	user.Email = "akdilsiz2@tecpor.com"

	err = db.Insert(new(TestUser), user, "id")
	assert.Nil(t, err)

	user2 := new(TestUser)
	user2.NodeID = node.ID
	user2.Username = "akdilsiz2"
	user2.PasswordDigest = "asdasdasd"
	user2.Email = "akdilsiz2@tecpor.com"

	err = db.Insert(new(TestUser), user2, "id")
	assert.NotNil(t, err)
	errs, err = ValidateConstraint(err, user2)
	assert.Equal(t, errs["username"], "has been already taken")

	user3 := new(TestUser)
	user3.NodeID = node.ID
	user3.Username = "akdilsiz2"
	user3.Email = "akdilsiz2@tecpor.com"

	err = db.Insert(new(TestUser), user3, "id")
	assert.NotNil(t, err)

	err = DropDB(db)
}
