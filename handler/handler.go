package handler

import (
	"encoding/json"
	"github.com/hayashiki/tarsier-integration/config"
	"github.com/hayashiki/tarsier-integration/pubsub"
	"github.com/hayashiki/tarsier-integration/repository"
	"github.com/hayashiki/tarsier-integration/slackauth"
	"log"
	"net/http"
)

type Handler struct {
	slackAuthSvc slackauth.Auth
	TeamRepo     repository.TeamRepository
	pubsubCli   pubsub.Client
}

func NewHandler() *Handler {
	conf := config.NewConfigFromEnv()
	teamRepo := repository.NewTeamRepository(repository.NewDSClient(conf.GCPProject))
	slackAuth := slackauth.NewAuth(
		conf.SlackClientID,
		conf.SlackSecretID,
		conf.SlackRedirectURL,
		slackauth.DefaultSlackBaseURL)
	pubsubCli, err := pubsub.NewClient("print_text", conf.GCPProject)
	if err != nil {
		log.Printf("fail to init pubsub, panic err=%v", err)
		panic(err)
	}

	return &Handler{
		slackAuthSvc: slackAuth,
		TeamRepo:     teamRepo,
		pubsubCli: pubsubCli,
	}
}

// Serve data as JSON as response
func jsonResponse(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
