package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/tarsier-integration/model"
)

//go:generate mockgen -source ./team_repository.go -destination ./mock/mock_team_repository.go
type TeamRepository interface {
	Put(team *model.Team) error
	GetByID(id string) (*model.Team, error)
}

type teamRepository struct {
	dsClient *datastore.Client
}

func (r teamRepository) GetByID(id string) (*model.Team, error) {
	ctx := context.Background()
	k := datastore.NameKey(model.TeamKind, id, nil)
	var dst model.Team
	err := r.dsClient.Get(ctx, k, &dst)
	return &dst, err
}

func (r teamRepository) Put(team *model.Team) error {
	ctx := context.Background()
	k := datastore.NameKey(model.TeamKind, team.ID, nil)
	_, err := r.dsClient.Put(ctx, k, team)
	if err != nil {
		return err
	}
	return nil
}

func NewTeamRepository(client *datastore.Client) TeamRepository {
	return &teamRepository{dsClient: client}
}

func NewDSClient(projectID string) *datastore.Client {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return client
}
