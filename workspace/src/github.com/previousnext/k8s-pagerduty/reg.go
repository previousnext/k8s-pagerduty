package main

import "regexp"

// Helper function to handle including and excluding from Pagerduty.
func reg(include, exclude, name string) (bool, error) {
	included, err := regexp.MatchString(include, name)
	if err != nil {
		return false, err
	}

	excluded, err := regexp.MatchString(exclude, name)
	if err != nil {
		return false, err
	}

	// If this is not meant to be on the list, stop it here.
	if excluded {
		return false, nil
	}

	return included, nil
}
