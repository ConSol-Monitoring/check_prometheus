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
func Query(address *url.URL, query, warning, critical, alias, search, replace, emptyQueryMessage string, emptyQueryStatus check_x.State) (err error) {
	warn, err := check_x.NewThreshold(warning)
	if err != nil {
		return err
	}

	crit, err := check_x.NewThreshold(critical)
	if err != nil {
		return err
	}
	var re *regexp.Regexp
	if search != "" {
		re, err = regexp.Compile(search)
		if err != nil {
			return err
		}
	}

	apiClient, err := helper.NewAPIClientV1(address)
	if err != nil {
		return err
	}

	result, _, err := apiClient.Query(context.TODO(), query, time.Now())
	if err != nil {
		return err
	}

	switch result.Type() {
	case model.ValScalar:
		scalar := result.(*model.Scalar)
		scalarValue := float64(scalar.Value)
		helper.CheckTimestampFreshness(scalar.Timestamp)

		check_x.NewPerformanceData(replaceLabel("scalar", re, replace), scalarValue).Warn(warn).Crit(crit)
		state := check_x.Evaluator{Warning: warn, Critical: crit}.Evaluate(scalarValue)

		resultAsString := strconv.FormatFloat(scalarValue, 'f', -1, 64)
		if alias == "" {
			check_x.Exit(state, fmt.Sprintf("Query: '%s' returned: '%s'", query, resultAsString))
		} else {
			check_x.Exit(state, fmt.Sprintf("Alias: '%s' returned: '%s'", alias, resultAsString))
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
			check_x.Exit(emptyQueryStatus, output)
		}
		for _, sample := range vector {
			helper.CheckTimestampFreshness(sample.Timestamp)

			sampleValue := float64(sample.Value)
			check_x.NewPerformanceData(replaceLabel(model.LabelSet(sample.Metric).String(), re, replace), sampleValue).Warn(warn).Crit(crit)
			states = append(states, check_x.Evaluator{Warning: warn, Critical: crit}.Evaluate(sampleValue))
			output += expandAlias(alias, sample.Metric, sampleValue)
		}

		return evalStates(states, output, query)
	case model.ValMatrix:
		matrix := result.(model.Matrix)
		states := check_x.States{}
		for _, sampleStream := range matrix {
			for _, value := range sampleStream.Values {
				helper.CheckTimestampFreshness(value.Timestamp)
				states = append(states, check_x.Evaluator{Warning: warn, Critical: crit}.Evaluate(float64(value.Value)))
			}
		}

		return evalStates(states, alias, query)
	default:
		check_x.Exit(check_x.Unknown, fmt.Sprintf("The query did not return a supported type(scalar, vector, matrix), instead: '%s'. Query: '%s'", result.Type().String(), query))

		return nil
	}

	return err
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

func evalStates(states check_x.States, alias, query string) error {
	state, err := states.GetWorst()
	if err != nil {
		return err
	}
	if alias == "" {
		check_x.Exit(*state, fmt.Sprintf("Query: '%s'", query))
	} else {
		check_x.Exit(*state, alias)
	}

	return nil
}
