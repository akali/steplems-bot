package config

import (
	"github.com/akali/steplems-bot/app/logger"
	"os"
	"strconv"
)

const (
	// BotAPITokenEnv is env var key for BotAPIToken.
	BotAPITokenEnv = "BOT_API_TOKEN"
	// UpdateTimeoutEnv is env var key for UpdateTimeout.
	UpdateTimeoutEnv = "UPDATE_TIMEOUT"
	// MongoConnectionStringEnv is env var key for connection string to mongodb.
	MongoConnectionStringEnv = "MONGO_CONNECTION_STRING"
	// MongoDatabaseNameEnv is env var key for database name in mongodb.
	MongoDatabaseNameEnv = "MONGO_DATABASE_NAME"
)

var (
	log = logger.Factory.NewLogger("config")

	// BotAPIToken is used to connect to bot api.
	BotAPIToken = os.Getenv(BotAPITokenEnv)
	// UpdateTimeout is a duration in seconds for bot api update chan timout.
	UpdateTimeout = 60
	// MongoConnectionString is connection string to mongodb.
	MongoConnectionString = os.Getenv(MongoConnectionStringEnv)
	// MongoDatabaseName is database name in mongodb.
	MongoDatabaseName = os.Getenv(MongoDatabaseNameEnv)
)

func init() {
	ut, err := strconv.ParseInt(os.Getenv(UpdateTimeoutEnv), 10, 32)
	if err != nil {
		log.Warn.Println("update timeout has not been set properly, using the default")
	} else {
		UpdateTimeout = int(ut)
	}

	log.Info.Println("connecting with connection string: ", MongoConnectionString)
	log.Info.Println("connecting with database name: ", MongoDatabaseName)
}
