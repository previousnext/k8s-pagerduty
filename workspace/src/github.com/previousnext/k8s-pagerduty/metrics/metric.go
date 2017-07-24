package metrics

import "fmt"

// Store contains metric data over 3 data points.
type Store struct {
	first  int
	second int
	third  int
}

// Add is used to add a new data point to the store.
func (m *Store) Add(val int) error {
	// Ensure that this value is within the range of 0 to 100.
	if val > 100 {
		return fmt.Errorf("value is greater than 100")
	}

	if val < 0 {
		return fmt.Errorf("value is less than 0")
	}

	// Update our storage.
	m.third = m.second
	m.second = m.first
	m.first = val

	return nil
}

// Avg is used to provide an average value over a dataset.
func (m *Store) Avg() (int, error) {
	if m.first == 0 {
		return m.first, fmt.Errorf("first data point does not contain data (or zero value)")
	}

	if m.second == 0 {
		return m.first, fmt.Errorf("second data point does not contain data (or zero value)")
	}

	if m.third == 0 {
		return m.first, fmt.Errorf("third data point does not contain data (or zero value)")
	}

	return (m.first + m.second + m.third) / 3, nil
}
