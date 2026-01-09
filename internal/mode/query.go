package mode

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/consol-monitoring/check_prometheus/internal/helper"

	"github.com/consol-monitoring/check_x"
	"github.com/prometheus/common/model"
)

// Query allows the user to test data in the prometheus server
func Query(ctx context.Context, address *url.URL, query, warning, critical, alias, search, replace, emptyQueryMessage string, emptyQueryStatus check_x.State, collection *check_x.PerformanceDataCollection) (check_x.State, string, error) {
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

	var re *regexp.Regexp
	if search != "" {
		re, err = regexp.Compile(search)
		if err != nil {
			return check_x.Unknown, fmt.Sprintf("Error creating regex from '%s' : %s", search, err.Error()), err
		}
	}

	apiClient, err := helper.NewAPIClientV1(address)
	if err != nil {
		return check_x.Unknown, fmt.Sprintf("Error creating apiClient: %s", err.Error()), err
	}

	result, _, err := apiClient.Query(ctx, query, time.Now())
	if err != nil {
		return check_x.Unknown, fmt.Sprintf("Error when querying: %s", err.Error()), err
	}

	switch result.Type() {
	case model.ValScalar:
		scalar := result.(*model.Scalar)
		scalarValue := float64(scalar.Value)
		if err := helper.CheckTimestampFreshness(scalar.Timestamp); err != nil {
			return check_x.Unknown, fmt.Sprintf("Error when checking scalar timestamp freshness: %s", err.Error()), err
		}

		collection.AddPerformanceDataFloat64(replaceLabel("scalar", re, replace), scalarValue)
		collection.Warn("scalar", warnThreshold)
		collection.Crit("scalar", critThreshold)
		state := check_x.Evaluator{Warning: warnThreshold, Critical: critThreshold}.Evaluate(scalarValue)

		resultAsString := strconv.FormatFloat(scalarValue, 'f', -1, 64)
		if alias == "" {
			return state, fmt.Sprintf("Query: '%s' returned: '%s'", query, resultAsString), nil
		} else {
			return state, fmt.Sprintf("Alias: '%s' returned: '%s'", alias, resultAsString), nil
		}
	case model.ValVector:
		vector := result.(model.Vector)
		states := check_x.States{}
		var output string
		if len(vector) == 0 && emptyQueryMessage != "" {
			output = emptyQueryMessage
		} else if len(vector) == 0 {
			output = fmt.Sprintf("Query '%s' returned no data.", query)
		}
		if output != "" {
			return emptyQueryStatus, output, nil
		}
		for _, sample := range vector {
			if err := helper.CheckTimestampFreshness(sample.Timestamp); err != nil {
				return check_x.Unknown, fmt.Sprintf("Error when checking sample timestamp freshness: %s", err.Error()), err
			}

			sampleValue := float64(sample.Value)
			label := replaceLabel(model.LabelSet(sample.Metric).String(), re, replace)
			collection.AddPerformanceDataFloat64(label, sampleValue)
			collection.Warn(label, warnThreshold)
			collection.Crit(label, critThreshold)
			states = append(states, check_x.Evaluator{Warning: warnThreshold, Critical: critThreshold}.Evaluate(sampleValue))
			output += expandAlias(alias, sample.Metric, sampleValue)
		}

		return evalStates(states, output, query)
	case model.ValMatrix:
		matrix := result.(model.Matrix)
		states := check_x.States{}
		for _, sampleStream := range matrix {
			for _, value := range sampleStream.Values {
				if err := helper.CheckTimestampFreshness(value.Timestamp); err != nil {
					return check_x.Unknown, fmt.Sprintf("Error when checking value timestamp freshness: %s", err.Error()), err
				}
				states = append(states, check_x.Evaluator{Warning: warnThreshold, Critical: critThreshold}.Evaluate(float64(value.Value)))
			}
		}

		return evalStates(states, alias, query)
	default:

		err = fmt.Errorf("query did not return a supported type(scalar, vector, matrix), instead: '%s'. Query: '%s'", result.Type().String(), query)

		return check_x.Unknown, fmt.Sprintf("Error when querying prometheus: %s", err.Error()), err
	}
}

func expandAlias(alias string, labels model.Metric, value float64) string {
	_, err := template.New("Output").Parse(alias)
	var output string
	if err != nil {
		output = alias
	} else {
		labelMap := make(map[string]string)
		for label, value := range labels {
			var l = fmt.Sprintf("%v", label)
			var v = fmt.Sprintf("%v", value)
			labelMap[l] = v
		}
		labelMap["xvalue"] = fmt.Sprintf("%v", value)
		var rendered bytes.Buffer
		output = rendered.String()
	}

	return output
}

func replaceLabel(label string, re *regexp.Regexp, replace string) string {
	if re != nil {
		label = re.ReplaceAllString(label, replace)
	}

	return label
}

func evalStates(states check_x.States, alias, query string) (check_x.State, string, error) {
	state, err := states.GetWorst()

	if err != nil {
		return *state, fmt.Sprintf("Error when picking the worst state out of many states: %s", err.Error()), err
	}

	if alias == "" {
		return *state, fmt.Sprintf("Query: '%s'", query), nil
	} else {
		return *state, alias, nil
	}
}
