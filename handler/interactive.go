package handler

import (
	"encoding/json"
	"fmt"
	"github.com/hayashiki/tarsier-integration/usecase"
	"github.com/slack-go/slack"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func (h *Handler) HandleSlackInteractive(w http.ResponseWriter, r *http.Request) error {

	fmt.Printf("Receive HandleSlackInteractive")
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to read request body: %v", err)
		return err
	}

	msg, err := interactionCallbackParse(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to read request body: %v", err)
		return err
	}

	switch msg.Type {
	case slack.InteractionTypeMessageAction:
		payload := &usecase.IssueDialogPayload{
			TriggerID:   msg.TriggerID,
			ChannelID:   msg.Channel.ID,
			ChannelName: msg.Channel.Name,
			TeamID:      msg.Team.ID,
			TeamDomain:  msg.Team.Domain,
			MessageTs:   msg.MessageTs,
			UserID:      msg.User.ID,
			MessageText: msg.Message.Text,
			File: usecase.File{
				Title:              msg.Message.Msg.Files[0].Title,
				Name:               msg.Message.Msg.Files[0].Name,
				URLPrivateDownload: msg.Message.Msg.Files[0].URLPrivateDownload,
			},
		}

		if err := h.pubsubCli.PublishCreateTopic(payload); err != nil {
			// TODO: http error
			return err
		}
	default:
		log.Println("Nothing to do")
	}

	w.WriteHeader(http.StatusOK)

	return err
}

func interactionCallbackParse(reqBody []byte) (*slack.InteractionCallback, error) {
	var req slack.InteractionCallback
	jsonStr, err := url.QueryUnescape(string(reqBody)[8:])
	err = json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		return nil, fmt.Errorf("error parsing interaction callback: Body: %s | Err: %s", reqBody, err)
	}
	return &req, nil
}
