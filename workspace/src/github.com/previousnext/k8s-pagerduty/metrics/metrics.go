package metrics

import (
	"fmt"
	"time"
)

// List is a collection of metrics.
type List struct {
	Items map[string]*Metrics
}

// Metrics are a collection container metrics.
type Metrics struct {
	LastUpdate time.Time
	CPU        *Store
	Memory     *Store
}

// New returns a new metrics list with garage collection enabled.
func New() List {
	// Initialise our storage.
	list := List{
		Items: make(map[string]*Metrics),
	}

	// Startup a GC background task.
	go list.garbageCollection()

	return list
}

// Helper function to run garbage collection.
func (m List) garbageCollection() {
	go func() {
		limiter := time.Tick(time.Minute)

		for {
			<-limiter

			for key, metric := range m.Items {
				if time.Since(metric.LastUpdate).Hours() > 1 {
					delete(m.Items, key)
				}
			}
		}
	}()
}

// Add stores new data points to the metric store object.
func (m List) Add(key string, cpu, memory int) error {
	// Ensure that this object exists in our store first.
	metric := getMetric(m, key)

	err := metric.CPU.Add(cpu)
	if err != nil {
		return fmt.Errorf("%s: failed to add cpu value: %s", key, err)
	}

	err = metric.Memory.Add(memory)
	if err != nil {
		return fmt.Errorf("%s: failed to add memory value: %s", key, err)
	}

	metric.LastUpdate = time.Now()

	m.Items[key] = metric

	return nil
}

// Helper function to return an initialized metrics object.
func getMetric(list List, key string) *Metrics {
	if val, ok := list.Items[key]; ok {
		return val
	}

	return &Metrics{
		LastUpdate: time.Now(),
		CPU:        &Store{},
		Memory:     &Store{},
	}
}

// AvgCPU returns the average CPU for a container.
func (m List) AvgCPU(key string) (int, error) {
	if val, ok := m.Items[key]; ok {
		return val.CPU.Avg()
	}

	return 0, fmt.Errorf("cannot find metric with key: %s", key)
}

// AvgMemory returns the average memory for a container.
func (m List) AvgMemory(key string) (int, error) {
	if val, ok := m.Items[key]; ok {
		return val.Memory.Avg()
	}

	return 0, fmt.Errorf("cannot find metric with key: %s", key)
}
