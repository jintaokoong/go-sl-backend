package environments

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"rmrf-slash.com/backend/configurations/logger"
)

type Configuration struct {
	MongoUsername   string `env:"MONGO_USERNAME"`
	MongoPassword   string `env:"MONGO_PASSWORD"`
	MongoConnection string `env:"MONGO_CONN"`
	FirebaseCreds   string `env:"FIREBASE_CREDS"`
}

func GetVariables() Configuration {
	// load env vars
	godotenv.Load()
	// bind struct
	cfg := Configuration{}
	if err := env.Parse(&cfg); err != nil {
		logger.GetInstance().Println("Failed to parse env vars")
	}
	return cfg
}
