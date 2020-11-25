package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GCPProject          string `envconfig:"GCP_PROJECT" required:"true"`
	Slack
}

type Slack struct {
	SlackClientID        string `envconfig:"SLACK_CLIENT_ID" required:"true"`
	SlackSecretID        string `envconfig:"SLACK_SECRET_ID" required:"true"`
	SlackRedirectURL     string `envconfig:"SLACK_REDIRECT_URL" required:"true"`
}

func NewConfigFromEnv() Config {
	env := Config{}
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalln("[ERROR] Can not read environment variables ", err)
		panic(err)
	}
	return env
}
