package main

import (
	"github.com/PagerDuty/go-pagerduty"
)

func pagerdutyEvent(action, serviceKey, indicentKey, description string) error {
	_, err := pagerduty.CreateEvent(pagerduty.Event{
		Type:        action,
		ServiceKey:  serviceKey,
		IncidentKey: indicentKey,
		Description: description,
	})
	return err
}
