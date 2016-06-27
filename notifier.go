package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/steemwatch/status/checks"

	"github.com/pkg/errors"
	"github.com/steemwatch/status/notifiers"
	"gopkg.in/yaml.v2"
)

func startNotifier() (chan<- checks.CheckSummary, error) {
	content, err := ioutil.ReadFile("./notifier_email.config.yml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "could not read email notifier config")
	}

	var config notifiers.EmailNotifierConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal email notifier config")
	}

	notifier, err := notifiers.NewEmailNotifier(&config)
	if err != nil {
		return nil, err
	}

	statusCh := make(chan checks.CheckSummary)
	go func() {
		for {
			summary := <-statusCh
			log.Println("Status change detected, sending notifications")
			if err := notifier.DispatchNotification(&summary); err != nil {
				log.Printf("Failed to send notification: %+v\n", err)
			}
		}
	}()
	return statusCh, nil
}
