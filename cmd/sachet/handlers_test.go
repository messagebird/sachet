package main

import (
	"testing"

	"github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/assert"
)

func Test_newAlertText(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		data template.Data
		exp  string
	}{
		{
			name: "empty",
			data: template.Data{},
			exp:  "Alert \n",
		},
		{
			name: "empty alerts",
			data: template.Data{
				Alerts: template.Alerts{
					template.Alert{},
				},
			},
			exp: " \n",
		},
		{
			name: "alerts labels",
			data: template.Data{
				Alerts: template.Alerts{
					template.Alert{
						Labels: map[string]string{
							"alertname":         "a",
							"instance":          "a",
							"exported_instance": "a",
						},
					},
				},
			},
			exp: " \nalertname= a\ninstance= a\nexported_instance= a",
		},
		{
			name: "common labels",
			data: template.Data{
				CommonLabels: template.KV{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			exp: "Alert \na | b | c",
		},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.exp, newAlertText(tc.data), tc.name)
	}
}
