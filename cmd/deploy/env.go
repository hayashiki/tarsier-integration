package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("./app.deploy.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := string(b)
	projectID := os.Args[1]
	slackSecretID := os.Args[2]
	slackRedirectURL := os.Args[3]
	slackAppID := os.Args[4]
	yaml := strings.Replace(lines, "##GCP_PROJECT", projectID, 1)
	yaml = strings.Replace(lines, "##SLACK_CLIENT_ID", slackSecretID, 1)
	yaml = strings.Replace(lines, "##SLACK_SECRET_ID", slackRedirectURL, 1)
	yaml = strings.Replace(lines, "##SLACK_REDIRECT_URL", slackAppID, 1)
	err = ioutil.WriteFile("./app.yaml", []byte(yaml), 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
