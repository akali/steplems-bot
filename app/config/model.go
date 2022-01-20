package config

import (
	"os"
	"strconv"

	"github.com/akali/steplems-bot/app/logger"
)

const (
	// BotAPITokenEnv is env var key for BotAPIToken.
	BotAPITokenEnv = "BOT_API_TOKEN"
	// UpdateTimeoutEnv is env var key for UpdateTimeout.
	UpdateTimeoutEnv = "UPDATE_TIMEOUT"

	// EnableMongoEnv is env var key for mongo enable check.
	EnableMongoEnv = "ENABLE_MONGO"

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
	// MongoDatabaseName is mongo enable flag.
	EnableMongo = false
	// MongoConnectionString is connection string to mongodb.
	MongoConnectionString = os.Getenv(MongoConnectionStringEnv)
	// MongoDatabaseName is database name in mongodb.
	MongoDatabaseName = os.Getenv(MongoDatabaseNameEnv)
)

// Init initializes env variables.
func Init() {
	ut, err := strconv.ParseInt(os.Getenv(UpdateTimeoutEnv), 10, 32)
	if err != nil {
		log.Warn.Println("update timeout has not been set properly, using the default")
	} else {
		UpdateTimeout = int(ut)
	}

	enableMongo, err := strconv.ParseBool(os.Getenv(EnableMongoEnv))
	if err != nil {
		log.Warn.PrintT("failed to read ENABLE_MONGO value, expected boolean, found {{value}}\n", os.Getenv(EnableMongoEnv))
	} else {
		EnableMongo = enableMongo
	}

	log.Info.Println("connecting with connection string: ", MongoConnectionString)
	log.Info.Println("connecting with database name: ", MongoDatabaseName)
}
