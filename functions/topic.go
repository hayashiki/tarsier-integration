package functions

import (
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hayashiki/tarsier-integration/repository"
	"github.com/hayashiki/tarsier-integration/usecase"
	"os"
)

func HandleCreateTopic(ctx context.Context, msg PubSubMessage) error {
	var payload usecase.IssueDialogPayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal data")
	}

	//h := handler.NewHandler()
	params := usecase.PrintTextParams{Payload: &payload}

	teamRepo := repository.NewTeamRepository(getDSClient(os.Getenv("GCP_PROJECT")))
	//h.TeamRepo
	// vupしてさしかえる
	uc := usecase.NewPrintText(teamRepo)

	err := uc.Do(params)

	if err != nil {
		return fmt.Errorf("failed to unmarshal data")
	}
	return nil
}

// TODO: replace NewDSClient
func getDSClient(projectID string) *datastore.Client {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return client
}
