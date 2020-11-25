package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/hayashiki/tarsier-integration/usecase"
)

type Client interface {
	Publish(serialized []byte) error
	PublishCreateTopic(payload *usecase.IssueDialogPayload) error
	//PublishUpdateTopic(serialized []byte) error
}

type client struct {
	ctx    context.Context
	topic  string
	client *pubsub.Client
}

func (c client) Publish(serialized []byte) error {
	topic := c.client.Topic(c.topic)
	_, err := topic.Publish(c.ctx, &pubsub.Message{Data: serialized}).Get(c.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c client) PublishCreateTopic(payload *usecase.IssueDialogPayload) error {
	serialized, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if err := c.Publish(serialized); err != nil {
		return err
	}
	return nil
}

func NewClient(topic, projectID string) (Client, error) {
	ctx := context.Background()
	cli, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	psCli := client{
		ctx:    ctx,
		topic:  topic,
		client: cli,
	}
	return psCli, err
}
