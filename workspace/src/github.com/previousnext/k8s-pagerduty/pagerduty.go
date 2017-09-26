package main

import (
	"log"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/cenkalti/backoff"
)

func pagerdutyEvent(action, serviceKey, indicentKey, description string) (*pagerduty.EventResponse, error) {
	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

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
			log.Printf("Failed to retrieve org members, retrying: %s", err)
			continue
		}

		ticker.Stop()
		break
	}

	return resp, err
}
