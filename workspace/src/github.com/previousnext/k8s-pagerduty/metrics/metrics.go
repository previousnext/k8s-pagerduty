package metrics

import (
	"fmt"
	"time"
)

type MetricsList struct {
	Items map[string]*Metrics
}

type Metrics struct {
	LastUpdate time.Time
	CPU        *Store
	Memory     *Store
}

func New() MetricsList {
	// Initialise our storage.
	list := MetricsList{
		Items: make(map[string]*Metrics),
	}

	// Startup a GC background task.
	go list.GarbageCollection()

	return list
}

func (m MetricsList) GarbageCollection() {
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

func (m MetricsList) Add(key string, cpu, memory int) error {
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

func getMetric(list MetricsList, key string) *Metrics {
	if val, ok := list.Items[key]; ok {
		return val
	}

	return &Metrics{
		LastUpdate: time.Now(),
		CPU:        &Store{},
		Memory:     &Store{},
	}
}

func (m MetricsList) AvgCPU(key string) (int, error) {
	if val, ok := m.Items[key]; ok {
		return val.CPU.Avg(), nil
	}

	return 0, fmt.Errorf("cannot find metric with key: %s", key)
}

func (m MetricsList) AvgMemory(key string) (int, error) {
	if val, ok := m.Items[key]; ok {
		return val.Memory.Avg(), nil
	}

	return 0, fmt.Errorf("cannot find metric with key: %s", key)
}
