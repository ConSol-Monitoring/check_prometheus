package mode

import (
	"testing"

	"github.com/prometheus/common/model"
)

func TestExpandAlias(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		labels   model.Metric
		value    float64
		expected string
	}{
		{
			name:     "simple template with label",
			alias:    "Hostname: {{.hostname}}",
			labels:   model.Metric{"hostname": "server01"},
			value:    42.0,
			expected: "Hostname: server01",
		},
		{
			name:     "template with multiple labels",
			alias:    "Host: {{.hostname}} Job: {{.job}}",
			labels:   model.Metric{"hostname": "server01", "job": "prometheus"},
			value:    1.0,
			expected: "Host: server01 Job: prometheus",
		},
		{
			name:     "template with xvalue",
			alias:    "Value: {{.xvalue}}",
			labels:   model.Metric{"hostname": "server01"},
			value:    123.45,
			expected: "Value: 123.45",
		},
		{
			name:     "template with if-else condition",
			alias:    "Status: {{if eq .xvalue \"1\"}}UP{{else}}DOWN{{end}}",
			labels:   model.Metric{"hostname": "server01"},
			value:    1.0,
			expected: "Status: UP",
		},
		{
			name:     "template with if-else condition false",
			alias:    "Status: {{if eq .xvalue \"1\"}}UP{{else}}DOWN{{end}}",
			labels:   model.Metric{"hostname": "server01"},
			value:    0.0,
			expected: "Status: DOWN",
		},
		{
			name:     "invalid template fallback",
			alias:    "{{invalid template",
			labels:   model.Metric{"hostname": "server01"},
			value:    42.0,
			expected: "{{invalid template",
		},
		{
			name:     "empty alias",
			alias:    "",
			labels:   model.Metric{"hostname": "server01"},
			value:    42.0,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandAlias(tt.alias, tt.labels, tt.value)
			if result != tt.expected {
				t.Errorf("expandAlias(%q, ...) = %q, want %q", tt.alias, result, tt.expected)
			}
		})
	}
}

