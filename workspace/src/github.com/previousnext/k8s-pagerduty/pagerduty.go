package main

import (
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
)

func pagerdutyEvent(action, serviceKey, indicentKey, description string) error {
	log.Infof("PagerDuty %s incident - %s", action, description)
	_, err := pagerduty.CreateEvent(pagerduty.Event{
		Type:        action,
		ServiceKey:  serviceKey,
		IncidentKey: indicentKey,
		Description: description,
	})
	return err
}
