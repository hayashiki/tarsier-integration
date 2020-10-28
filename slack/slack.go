package slack

import (
	"bytes"
	"github.com/slack-go/slack"
	"io"
	"log"
)

//go:generate mockgen -source ./slack.go -destination ./mock/mock_slack.go
type Slack interface {
	Upload(title, name, channel, ts string, r io.Reader) error
	Download(url string, b *bytes.Buffer) error
}

type client struct {
	cli *slack.Client
}

func NewClient(token string) Slack {
	cli := slack.New(token)

	return &client{
		cli: cli,
	}
}

func (c *client) Download(url string, b *bytes.Buffer) error {

	err := c.cli.GetFile(url, b)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *client) Upload(title, name, channel, ts string, r io.Reader) error {

	params := slack.FileUploadParameters{
		Title:  title,
		File:   name,
		Reader: r,
		//InitialComment: TODO: need??,
		Channels:        []string{channel},
		ThreadTimestamp: ts,
	}

	_, err := c.cli.UploadFile(params)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
