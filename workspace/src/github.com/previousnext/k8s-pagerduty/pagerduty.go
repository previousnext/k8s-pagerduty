package main

import (
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/cenkalti/backoff"
)

func pagerdutyEvent(action, serviceKey, indicentKey, description string) (*pagerduty.EventResponse, error) {
	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

	log.Infof("PagerDuty %s incident - %s", action, description)

	var (
		resp *pagerduty.EventResponse
		err  error
	)

	for range ticker.C {
		resp, err = pagerduty.CreateEvent(pagerduty.Event{
			Type:        action,
			ServiceKey:  serviceKey,
			IncidentKey: indicentKey,
			Description: description,
		})
		if err != nil {
			log.Infof("Failed to create PagerDuty event, retrying: %s", err)
			continue
		}

		ticker.Stop()
		break
	}

	return resp, err
}
