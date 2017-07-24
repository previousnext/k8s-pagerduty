package main

import (
	"fmt"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/previousnext/k8s-pagerduty/metrics"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cliIncludePod       = kingpin.Flag("include-pod", "A regex which allows the user to include alerts to a specific pod eg. prod-*").Default("prod-*").OverrideDefaultFromEnvar("FILTER_POD").String()
	cliExcludePod       = kingpin.Flag("exclude-pod", "A regex which allows the user to exclude alerts to a specific pod eg. prod-backup*").Default("prod-backup*").OverrideDefaultFromEnvar("FILTER_POD").String()
	cliIncludeContainer = kingpin.Flag("include-container", "A regex which allows the user to filter alerts to a specific pod eg. prod-*").Default("^app$").OverrideDefaultFromEnvar("FILTER_CONTAINER").String()
	cliEndpoint         = kingpin.Flag("docker", "The Docker API endpoint").Default("unix:///var/run/docker.sock").OverrideDefaultFromEnvar("DOCKER_API").String()
	cliInterval         = kingpin.Flag("interval", "How long to wait before checking CPU and Memory again").Default("120s").OverrideDefaultFromEnvar("INTERVAL").Duration()
	cliServiceKey       = kingpin.Flag("service", "PagerDuty service key used for triggers").Default("").OverrideDefaultFromEnvar("SERVICE_KEY").String()
	cliThresholdCPU     = kingpin.Flag("threshold-cpu", "Threshold to determine if an application is using too much CPU").Default("95").OverrideDefaultFromEnvar("THRESHOLD_CPU").Int()
	cliThresholdMemory  = kingpin.Flag("threshold-memory", "Threshold to determine if an application is using too much memory").OverrideDefaultFromEnvar("THRESHOLD_MEMORY").Default("95").Int()
)

func main() {
	kingpin.Parse()

	fmt.Println("Starting PagerDuty Notifier")

	d, err := docker.NewClient(*cliEndpoint)
	if err != nil {
		panic(err)
	}

	var (
		store   = metrics.New()
		limiter = time.Tick(*cliInterval)
	)

	for {
		<-limiter

		fmt.Println("Looking up containers")

		containers, err := d.ListContainers(docker.ListContainersOptions{})
		if err != nil {
			panic(err)
		}

		// Load the existing containers that are running on this host.
		for _, container := range containers {
			namespace, pod, name, err := filter(container.Labels)
			if err != nil {
				fmt.Println(container.ID, "does not fit criteria:", err)
				continue
			}

			// The following allows us to filter our pods down to a specific subset.
			//   eg. Only alert me for pods matching "prod-*" (who cares about dev/stg).
			matched, err := reg(*cliIncludePod, *cliExcludePod, pod)
			if err != nil {
				fmt.Println(pod, "does not compily with regex:", err)
				continue
			}

			if !matched {
				fmt.Println(pod, "does not match regex")
				continue
			}

			// The following allows us to filter our container down to a specific subset.
			//   eg. Only alert me for pods matching "app".
			matched, err = regexp.MatchString(*cliIncludeContainer, name)
			if err != nil {
				fmt.Println(name, "does not compily with regex:", err)
				continue
			}

			if !matched {
				fmt.Println(name, "does not match regex")
				continue
			}

			info, err := d.InspectContainer(container.ID)
			if err != nil {
				fmt.Println("unable to inspect container", container.ID, err)
				continue
			}

			v, err := stats(d, container.ID)
			if err != nil {
				fmt.Println("unable to get stats for", container.ID, err)
				continue
			}

			if info.HostConfig.CPUQuota > 0 && info.HostConfig.Memory > 0 {
				var (
					cpuDelta    = float64(v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage)
					systemDelta = float64(v.CPUStats.SystemCPUUsage - v.PreCPUStats.SystemCPUUsage)
					cpu         = int((cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0)
					mem         = int(float64(v.MemoryStats.Usage) / float64(info.HostConfig.Memory) * 100)
				)

				log.WithFields(log.Fields{
					"namespace": namespace,
					"pod":       pod,
					"name":      name,
					"cpu":       cpu,
					"memory":    mem,
				}).Info("storing the following cpu and memory data")

				key := fmt.Sprintf("%s/%s/%s", namespace, pod, name)

				// Add our latest data to this entry.
				store.Add(key, cpu, mem)

				avgCPU, err := store.AvgCPU(key)
				if err != nil {
					fmt.Println("unable to get average cpu metric:", err)
					continue
				}

				avgMemory, err := store.AvgMemory(key)
				if err != nil {
					fmt.Println("unable to get average memory metric:", err)
					continue
				}

				// Compare the current CPU metric and the previous memory metric vs threshold.
				if avgCPU > *cliThresholdCPU {
					err := pagerdutyEvent("trigger", *cliServiceKey, key, fmt.Sprintf("CPU load GREATER than %v (%v) for container %s", *cliThresholdCPU, cpu, key))
					if err != nil {
						panic(err)
					}
				} else {
					err := pagerdutyEvent("resolve", *cliServiceKey, key, fmt.Sprintf("CPU load LESS than %v (%v) for container %s", *cliThresholdCPU, cpu, key))
					if err != nil {
						panic(err)
					}
				}

				// Compare the current memory metric and the previous memory metric vs threshold.
				if avgMemory > *cliThresholdMemory {
					err := pagerdutyEvent("trigger", *cliServiceKey, key, fmt.Sprintf("Memory usage GREATER than %v (%v) for container %s", *cliThresholdMemory, cpu, key))
					if err != nil {
						panic(err)
					}
				} else {
					err := pagerdutyEvent("resolve", *cliServiceKey, key, fmt.Sprintf("Memory usage LESS than %v (%v) for container %s", *cliThresholdMemory, cpu, key))
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
