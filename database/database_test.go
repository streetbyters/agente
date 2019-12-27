package database

import (
	"github.com/akdilsiz/agente/model"
	"github.com/akdilsiz/agente/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var logger = utils.NewLogger("test")
var appPath, _ = os.Getwd()
var dirs = strings.SplitAfter(appPath, "agente")

func Test_NewDB(t *testing.T) {
	appPath = dirs[0]

	// Open postgres db connection
	config := &model.Config{
		Mode:   model.Test,
		DB:     model.Postgres,
		DBName: "agente_test",
		DBHost: "127.0.0.1",
		DBPort: 5432,
		DBUser: "agente",
		DBPass: "123456",
		DBSsl:  "disable",
	}

	_, err := NewDB(config)
	if err != nil {
		t.Fatal(err)
	}

	logger.LogInfo("Success open postgres db connection")

	// Failed postgres db connection if given invalid information
	config = &model.Config{
		Mode:   model.Test,
		DB:     model.Postgres,
		DBName: "agente-error",
		DBHost: "127.0.0.1",
		DBPort: 5432,
		DBUser: "agente-error",
		DBPass: "123456",
		DBSsl:  "disable",
	}

	_, err = NewDB(config)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed postgres db connection if given invalid " +
		"information")

	// Failed postgres db connection if given invalid port
	config = &model.Config{
		Mode:   model.Test,
		DB:     model.Postgres,
		DBName: "agente_test",
		DBHost: "127.0.0.4",
		DBPort: 5435,
		DBUser: "agente",
		DBPass: "123456",
		DBSsl:  "disable",
	}

	_, err = NewDB(config)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed postgres db connection if given invalid " +
		"port")

	//Failed db connection if unsupported database type
	config = &model.Config{
		Mode:   model.Test,
		DB:     model.Unknown,
		DBName: "agente-error",
		DBHost: "127.0.0.1",
		DBPort: 3303,
		DBUser: "agente-error",
		DBPass: "123456",
		DBSsl:  "false",
	}

	_, err = NewDB(config)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed db connection if unsupported database type")
}

func Test_InstallDB(t *testing.T) {
	appPath = dirs[0]

	// Install postgres db
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

	database, err := NewDB(config)
	if err != nil {
		t.Fatal(err)
	}
	database.Logger = logger
	database.Reset = true
	DropDB(database)

	err = InstallDB(database)
	if err != nil {
		t.Fatal(err)
	}
	logger.LogInfo("InstallDB Successfully postgres. If no migration was made before.")

	// New migration up for postgres db
	config = &model.Config{
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

	database, err = NewDB(config)
	if err != nil {
		t.Fatal(err)
	}
	database.Logger = logger
	DropDB(database)

	err = InstallDB(database)
	if err != nil {
		t.Fatal(err)
	}
	logger.LogInfo("New migration up for postgres db")

	data, err := ioutil.ReadFile(filepath.Join(config.DBPath, "sql", "postgres", "02.base_tables.down.sql"))
	if err != nil {
		t.Fatal(err)
	}

	res := database.Query(string(data))
	if res.Error != nil {
		t.Fatal(res.Error)
	}
	data, err = ioutil.ReadFile(filepath.Join(config.DBPath, "sql", "postgres", "03.create_user_passphrases.down.sql"))
	if err != nil {
		t.Fatal(err)
	}

	res = database.Query(string(data))
	if res.Error != nil {
		t.Fatal(err)
	}

	res = database.Query("delete from " + string(tMigration) + " where id > 0")

	err = InstallDB(database)
	if err != nil {
		t.Fatal(err)
	}

	result := database.Query("select * from " + string(tMigration))
	if result.Error != nil {
		t.Fatal(err)
	}

	logger.LogInfo("Successfully new migrations")
}
