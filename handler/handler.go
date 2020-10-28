package handler

import (
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"github.com/hayashiki/tarsier-integration/repository"
	"github.com/hayashiki/tarsier-integration/slackauth"
	"net/http"
	"os"
)

type Handler struct {
	slackAuthSvc slackauth.Auth
	teamRepo     repository.TeamRepository
}

func NewHandler() *Handler {
	teamRepo := repository.NewTeamRepository(getDSClient(os.Getenv("GCP_PROJECT")))
	slackAuth := slackauth.NewAuth(
		os.Getenv("SLACK_CLIENT_ID"),
		os.Getenv("SLACK_SECRET_ID"),
		os.Getenv("SLACK_REDIRECT_URL"),
		slackauth.DefaultSlackBaseURL)

	return &Handler{
		slackAuthSvc: slackAuth,
		teamRepo:     teamRepo,
	}
}

// Serve data as JSON as response
func jsonResponse(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func getDSClient(projectID string) *datastore.Client {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return client
}
