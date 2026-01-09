package mode

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/consol-monitoring/check_prometheus/internal/helper"

	"github.com/consol-monitoring/check_x"
)

const (
	// DefaultLabel is used if the given label is wrong
	DefaultLabel = "instance"
)

type targets struct {
	Status string `json:"status"`
	Data   struct {
		ActiveTargets []struct {
			DiscoveredLabels struct {
				Address     string `json:"__address__"`
				MetricsPath string `json:"__metrics_path__"`
				Scheme      string `json:"__scheme__"`
				Job         string `json:"job"`
			} `json:"discoveredLabels"`
			Labels     map[string]string `json:"labels"`
			ScrapeURL  string            `json:"scrapeUrl"`
			LastError  string            `json:"lastError"`
			LastScrape time.Time         `json:"lastScrape"`
			Health     string            `json:"health"`
		} `json:"activeTargets"`
	} `json:"data"`
}

func getTargets(address *url.URL) (*targets, error) {
	u, err := url.Parse(address.String())
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api/v1/targets")
	jsonBytes, err := helper.DoAPIRequest(u)
	if err != nil {
		return nil, err
	}
	var dat targets
	if err = json.Unmarshal(jsonBytes, &dat); err != nil {
		return nil, err
	}

	return &dat, nil
}

// TargetsHealth tests the health of the targets
func TargetsHealth(address *url.URL, label, warning, critical string, collection *check_x.PerformanceDataCollection) (check_x.State, string, error) {
	if address == nil {
		err := fmt.Errorf("address to query is null")
		return check_x.Unknown, fmt.Sprintf("Error: %s", err.Error()), err
	}

	if collection == nil {
		err := fmt.Errorf("collection to store perf data is null")
		return check_x.Unknown, fmt.Sprintf("Error: %s", err.Error()), err
	}

	warnThreshold, err := check_x.NewThreshold(warning)
	if err != nil {
		return check_x.Unknown, fmt.Sprintf("Error creating warningThreshold from '%s' : %s", warning, err.Error()), err
	}

	critThreshold, err := check_x.NewThreshold(critical)
	if err != nil {
		return check_x.Unknown, fmt.Sprintf("Error creating critThreshold from '%s' : %s", critical, err.Error()), err
	}

	targets, err := getTargets(address)
	if err != nil {
		return check_x.Unknown, fmt.Sprintf("Error getting targets out of address: %s : %s", address.String(), err.Error()), err
	}

	if (*targets).Status != "success" {
		err := fmt.Errorf("the API target return status was %s", (*targets).Status)
		return check_x.Unknown, err.Error(), err
	}
	msg := ""
	healthy := 0
	unhealthy := 0
	for _, target := range (*targets).Data.ActiveTargets {
		msg += fmt.Sprintf("Job: %s, Instance: %s, Health: %s, Last Error: %s\n", target.Labels["job"], target.Labels["instance"], target.Health, target.LastError)
		health := 0.0
		if target.Health != "up" {
			health = 1
			unhealthy += 1
		} else {
			healthy += 1
		}
		if val, ok := target.Labels[label]; ok {
			collection.AddPerformanceDataFloat64(val, health)
		} else {
			collection.AddPerformanceDataFloat64(target.Labels[DefaultLabel], health)
		}
	}
	var healthRate float64
	sumTargets := float64(len((*targets).Data.ActiveTargets))
	if sumTargets == 0 {
		healthRate = 0
	} else {
		healthRate = float64(healthy) / sumTargets
	}

	collection.AddPerformanceDataFloat64("health_rate", healthRate)
	collection.Warn("health_rate", warnThreshold)
	collection.Crit("health_rate", critThreshold)
	collection.Min("health_rate", 0)
	collection.Max("health_rate", 1)
	collection.AddPerformanceDataFloat64("targets", sumTargets)
	collection.Min("targets", 0)

	state := check_x.Evaluator{Warning: warnThreshold, Critical: critThreshold}.Evaluate(healthRate)

	return state, fmt.Sprintf("There are %d healthy and %d unhealthy targets", healthy, unhealthy), nil
}
