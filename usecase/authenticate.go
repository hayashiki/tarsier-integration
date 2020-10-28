package usecase

import (
	"fmt"
	"github.com/hayashiki/tarsier-integration/model"
	"github.com/hayashiki/tarsier-integration/repository"
	"github.com/hayashiki/tarsier-integration/slackauth"
	"time"
)

type AuthenticateSlack struct {
	slackAuthSvc slackauth.Auth
	teamRepo     repository.TeamRepository
}

type AuthenticateSlackParams struct {
	Code string
}

type AuthenticateSlackResp struct {
	AuthResp slackauth.Response
}

func NewSlackAuthenticate(slackAuthSvc slackauth.Auth, teamRepo repository.TeamRepository) *AuthenticateSlack {
	return &AuthenticateSlack{
		slackAuthSvc: slackAuthSvc,
		teamRepo:     teamRepo,
	}
}

func (a *AuthenticateSlack) Do(params AuthenticateSlackParams) (*AuthenticateSlackResp, error) {
	authResp, err := a.slackAuthSvc.ExchangeSlackAuthCodeForToken(params.Code)

	if err != nil {
		return nil, fmt.Errorf("failed to get slack access err=%w", err)
	}

	team := &model.Team{
		ID:        authResp.Team.ID,
		Name:      authResp.Team.Name,
		Token:     authResp.AccessToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = a.teamRepo.Put(team)
	if err != nil {
		return nil, fmt.Errorf("failed to save slack team info %v", err)
	}

	return &AuthenticateSlackResp{AuthResp: authResp}, nil
}
