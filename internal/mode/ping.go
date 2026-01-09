package mode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/consol-monitoring/check_prometheus/internal/helper"
	"github.com/consol-monitoring/check_x"
	"github.com/prometheus/common/model"
)

type buildInfo struct {
	Metric struct {
		Name      string `json:"__name__"`
		Branch    string `json:"branch"`
		Goversion string `json:"goversion"`
		Instance  string `json:"instance"`
		Job       string `json:"job"`
		Revision  string `json:"revision"`
		Version   string `json:"version"`
	} `json:"metric"`
	Value []interface{} `json:"value"`
}

// Ping will fetch build information from the prometheus server
func Ping(ctx context.Context, address *url.URL, collection *check_x.PerformanceDataCollection) (check_x.State, string, error) {
	if address == nil {
		err := fmt.Errorf("address to query is null")
		return check_x.Unknown, fmt.Sprintf("Error: %s", err.Error()), err
	}

	if collection == nil {
		err := fmt.Errorf("collection to store perf data is null")
		return check_x.Unknown, fmt.Sprintf("Error: %s", err.Error()), err
	}

	apiClient, err := helper.NewAPIClientV1(address)
	if err != nil {
		return check_x.Unknown, fmt.Sprintf("Error creating apiClient: %s", err.Error()), err
	}

	query := `prometheus_build_info{job="prometheus"}`
	startTime := time.Now()
	result, _, err := apiClient.Query(ctx, query, time.Now())
	endTime := time.Now()
	if err != nil {
		return check_x.Unknown, fmt.Sprintf("Error when querying API: %s", err.Error()), err
	}

	vector := result.(model.Vector)
	if len(vector) != 1 {
		err := fmt.Errorf("the query '%s' did not return a vector with a single entry", query)
		return check_x.Unknown, fmt.Sprintf("Query returned a single element: %s", err.Error()), err
	}

	sample := vector[0]
	if err := helper.CheckTimestampFreshness(sample.Timestamp); err != nil {
		return check_x.Unknown, fmt.Sprintf("Error when checking sample timestamp freshness: %s", err.Error()), err
	}
	jsonBytes, err := sample.MarshalJSON()
	if err != nil {
		return check_x.Unknown, fmt.Sprintf("Error when marshalling json of first sample in the vector: %s", err.Error()), err
	}
	var dat buildInfo
	if err := json.Unmarshal(jsonBytes, &dat); err != nil {
		return check_x.Unknown, fmt.Sprintf("Error when unmarshalling json of first sample in the vector: %s", err.Error()), err
	}
	collection.AddPerformanceDataFloat64("duration", endTime.Sub(startTime).Seconds())
	collection.Unit("duration", "s")
	collection.Min("duration", 0)

	return check_x.OK, fmt.Sprintf("Version: %s, Instance %s", dat.Metric.Version, dat.Metric.Instance), nil
}
