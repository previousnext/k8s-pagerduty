package metrics

import "fmt"

type Store struct {
	first  int
	second int
	third  int
}

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

func (m *Store) Avg() int {
	return (m.first + m.second + m.third) / 3
}
