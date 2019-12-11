package cmn

import (
	"github.com/akdilsiz/agente/model"
	"os"
	"path"
	"strings"
	"testing"
)

var logger = NewLogger("test")
var appPath, _ = os.Getwd()
var dirs = strings.SplitAfter(appPath, "agente")

func Test_NewDB(t *testing.T) {
	appPath = dirs[0]

	// Open sqlite connection
	config := &model.Config{
		DBPath: appPath,
		Mode:   model.Test,
		DB:     model.SQLite,
		DBName: "agenteTest.db",
	}

	_, err := NewDB(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	logger.LogInfo("Success open sqlite connection")

	// Failed sqlite connection if file permission error
	config = &model.Config{
		DBPath: "/root",
		Mode:   model.Test,
		DB:     model.SQLite,
		DBName: "agenteTest.db",
	}

	_, err = NewDB(config, logger)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed sqlite connection if file permission " +
		"error")

	// Open postgres db connection
	config = &model.Config{
		Mode:   model.Test,
		DB:     model.Postgres,
		DBName: "agente",
		DBHost: "127.0.0.1",
		DBPort: 5432,
		DBUser: "agente",
		DBPass: "123456",
		DBSsl:  "disable",
	}

	_, err = NewDB(config, logger)
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

	_, err = NewDB(config, logger)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed postgres db connection if given invalid " +
		"information")

	// Failed postgres db connection if given invalid port
	config = &model.Config{
		Mode:   model.Test,
		DB:     model.Postgres,
		DBName: "agente",
		DBHost: "127.0.0.4",
		DBPort: 5435,
		DBUser: "agente",
		DBPass: "123456",
		DBSsl:  "disable",
	}

	_, err = NewDB(config, logger)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed postgres db connection if given invalid " +
		"port")

	// Open mysql db connection
	config = &model.Config{
		Mode:   model.Test,
		DB:     model.Mysql,
		DBName: "agente",
		DBHost: "127.0.0.1",
		DBPort: 3306,
		DBUser: "agente",
		DBPass: "123456",
		DBSsl:  "false",
	}

	_, err = NewDB(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	logger.LogInfo("Success open mysql db connection")

	// Failed mysql db connection if given invalid information
	config = &model.Config{
		Mode:   model.Test,
		DB:     model.Mysql,
		DBName: "agente-error",
		DBHost: "127.0.0.1",
		DBPort: 3306,
		DBUser: "agente-error",
		DBPass: "123456",
		DBSsl:  "false",
	}

	_, err = NewDB(config, logger)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed mysql db connection if given invalid " +
		"information")

	//Failed mysql db connection if given invalid port
	config = &model.Config{
		Mode:   model.Test,
		DB:     model.Mysql,
		DBName: "agente-error",
		DBHost: "127.0.0.1",
		DBPort: 3303,
		DBUser: "agente-error",
		DBPass: "123456",
		DBSsl:  "false",
	}

	_, err = NewDB(config, logger)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed mysql db connection if given invalid " +
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

	_, err = NewDB(config, logger)
	if err == nil {
		t.Fatal(err)
	}

	logger.LogInfo("Failed db connection if unsupported database type")
}

func Test_InstallDB(t *testing.T) {
	appPath = dirs[0]

	// Open sqlite connection
	config := &model.Config{
		DBPath: appPath,
		Mode:   model.Test,
		DB:     model.SQLite,
		DBName: "agenteTest.db",
	}

	database, err := NewDB(config, logger)
	if err != nil {
		t.Fatal(err)
	}
	DropDB(database)

	err = InstallDB(database)
	if err != nil {
		t.Fatal(err)
	}
	logger.LogInfo("Success install sqlite. If no migration was made before.")

	// Install postgres db
	config = &model.Config{
		DBPath: path.Join(appPath, "sql"),
		Mode:   model.Test,
		DB:     model.Postgres,
		DBName: "agente",
		DBHost: "127.0.0.1",
		DBPort: 5432,
		DBUser: "agente",
		DBPass: "123456",
		DBSsl:  "disable",
	}

	database, err = NewDB(config, logger)
	if err != nil {
		t.Fatal(err)
	}
	DropDB(database)

	err = InstallDB(database)
	if err != nil {
		t.Fatal(err)
	}
	logger.LogInfo("InstallDB Successfully postgres. If no migration was made before.")
}
