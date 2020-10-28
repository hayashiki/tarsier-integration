package usecase

import (
	"bytes"
	"fmt"
	"github.com/hayashiki/tarsier-integration/repository"
	"github.com/hayashiki/tarsier-integration/slack"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"path/filepath"

	"github.com/hayashiki/tarsier"

	"log"
	"os"
)

type IssueDialogPayload struct {
	TriggerID   string
	ChannelID   string
	ChannelName string
	TeamID      string
	TeamDomain  string
	MessageTs   string
	UserID      string
	MessageText string
	File        File
	Files       []File
}

type File struct {
	Title              string
	Name               string
	URLPrivateDownload string
}

type PrintText interface {
	Do(params PrintTextParams) error
}

type PrintTextParams struct {
	Payload *IssueDialogPayload
}

func NewPrintText(
	teamRepository repository.TeamRepository,
) PrintText {
	return &printText{
		teamRepo: teamRepository,
	}
}

type printText struct {
	teamRepo repository.TeamRepository
}

func (uc *printText) Do(params PrintTextParams) error {

	team, err := uc.teamRepo.GetByID(params.Payload.TeamID)

	if err != nil {
		log.Printf("GetByID err %v", err)
		return err
	}

	slackSvc := slack.NewClient(team.Token)
	// TODO: fileNameは拡張子で制御を入れる
	fileName := filepath.Base(params.Payload.File.URLPrivateDownload)
	b := bytes.Buffer{}
	err = slackSvc.Download(params.Payload.File.URLPrivateDownload, &b)

	if err != nil {
		log.Printf("failed to get a slack file err=%v", err)
		return err
	}
	img, _, err := image.Decode(&b)
	if err != nil {
		log.Printf("failed to get a slack file, err=%v", err)
		return err
	}

	outImg, err := tarsier.Print(img, tarsier.DefaultOverlayText)
	file, err := os.Create(fmt.Sprintf("/tmp/%s", fileName)) // os.Create(fmt.Sprintf("/tmp/%s", fileName), os.O_TRUNC|os.O_RDWR, 0755)
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	err = png.Encode(file, outImg)
	if err != nil {
		log.Println(err)
		return err
	}

	err = slackSvc.Upload(fileName, fmt.Sprintf("/tmp/%s", fileName), params.Payload.ChannelName, params.Payload.MessageTs, file)
	if err != nil {
		log.Printf("failed to upload a slack file err=%v", err)
		return err
	}

	return nil
}
