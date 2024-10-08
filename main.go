// Post build status results to Slack.

package main

import (
	"context"
	"flag"
	"log"
	"io/ioutil"

	"github.com/gijsentius/cloud-builders-slackbot/slackbot"
)

var (
	buildId     = flag.String("build", "", "Id of monitored Build")
	messageFile  = flag.String("file", "", "file with message to send")
	webhook     = flag.String("webhook", "", "Slack webhook URL")
	project     = flag.String("project", "unknown", "Project name being built")
	mode        = flag.String("mode", "trigger", "Mode the builder runs in")
	copyName    = flag.Bool("copy-name", false, "Copy name of slackbot's build step from monitored build to watcher build")
	copyTags    = flag.Bool("copy-tags", false, "Copy tags from monitored build to watcher build")
	copyTimeout = flag.Bool("copy-timeout", false, "Copy timeout from monitored build to watcher build")
)

func main() {
	log.Print("Starting slackbot")
	flag.Parse()
	ctx := context.Background()

	if *webhook == "" {
		log.Fatalf("Slack webhook must be provided.")
	}

	if *buildId == "" {
		log.Fatalf("Build ID must be provided.")
	}
	if *messageFile == "" {
		log.Fatalf("Message must be provided.")
	}

	if *mode == "" {
		log.Fatalf("Mode must be provided.")
	}

	if *mode != "trigger" && *mode != "monitor" {
		log.Fatalf("Mode must be one of: trigger, monitor.")
	}

	projectId, err := slackbot.GetProject()
	if err != nil {
		log.Fatalf("Failed to get project ID: %v", err)
	}

	messageBytes, err := ioutil.ReadFile(messageFile)
    if err != nil {
        log.Fatalf("Failed to parse message from file: %v", err)
    }
	message := string(messageBytes)

	if *mode == "trigger" {
		// Trigger another build to run the monitor.
		log.Printf("Starting trigger mode for build %s", *buildId)
		slackbot.Trigger(ctx, projectId, *buildId, *webhook, *project, *copyName, *copyTags, *copyTimeout, *message)
		return
	}

	if *mode == "monitor" {
		// Monitor the other build until completion.
		log.Printf("Starting monitor mode for build %s", *buildId)
		slackbot.Monitor(ctx, projectId, *buildId, *webhook, *project, *message)
		return
	}
}
