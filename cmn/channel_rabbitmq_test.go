package cmn

import (
	"errors"
	"github.com/forgolang/agente/model"
	"github.com/forgolang/agente/utils"
	"os"
	"testing"
)

var logger = utils.NewLogger("test")
var appPath, _ = os.Getwd()

func Test_NewRabbitMq(t *testing.T) {
	config := &model.Config{
		Path:         appPath,
		Mode:         model.Test,
		RabbitMq:     false,
		RabbitMqHost: "127.0.0.1",
		RabbitMqPort: 5672,
		RabbitMqUser: "local",
		RabbitMqPass: "local",
		ChannelName:  "agente_test",
		Versioning:   false,
		Scheduler:    "JobRunner",
	}
	app := NewApp(config, logger)

	channelRabbitMq := NewRabbitMq(app)
	if app != channelRabbitMq.App {
		t.Fatal(errors.New("NewRabbitMq error"))
	}
}
