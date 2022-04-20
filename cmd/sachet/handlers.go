package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/prometheus/alertmanager/template"

	"github.com/messagebird/sachet"
)

type handlers struct{}

func newAlertText(data template.Data) string {
	if len(data.Alerts) > 1 {
		labelAlerts := map[string]template.Alerts{
			"Firing":   data.Alerts.Firing(),
			"Resolved": data.Alerts.Resolved(),
		}
		text := ""
		for label, alerts := range labelAlerts {
			if len(alerts) > 0 {
				text += label + ": \n"
				for _, alert := range alerts {
					text += alert.Labels["alertname"] + " @" + alert.Labels["instance"]
					if len(alert.Labels["exported_instance"]) > 0 {
						text += " (" + alert.Labels["exported_instance"] + ")"
					}
					text += "\n"
				}
			}
		}
		return text
	}

	if len(data.Alerts) == 1 {
		alert := data.Alerts[0]
		tuples := []string{}
		for k, v := range alert.Labels {
			tuples = append(tuples, k+"= "+v)
		}
		sort.Strings(tuples)
		return strings.ToUpper(data.Status) + " \n" + strings.Join(tuples, "\n")
	}

	return "Alert \n" + strings.Join(data.CommonLabels.Values(), " | ")
}

func (h handlers) Alert(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// https://godoc.org/github.com/prometheus/alertmanager/template#Data
	data := template.Data{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errorHandler(w, http.StatusBadRequest, err, "?")
		return
	}

	receiverConf := receiverConfByReceiver(data.Receiver)
	if receiverConf == nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("Receiver missing: %s", data.Receiver), "?")
		return
	}
	provider, err := providerByName(receiverConf.Provider)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err, receiverConf.Provider)
		return
	}

	var text string
	if receiverConf.Text != "" {
		text, err = tmpl.ExecuteTextString(receiverConf.Text, data)
		if err != nil {
			errorHandler(w, http.StatusInternalServerError, err, receiverConf.Provider)
			return
		}
	} else {
		text = newAlertText(data)
	}

	message := sachet.Message{
		To:   receiverConf.To,
		From: receiverConf.From,
		Type: receiverConf.Type,
		Text: text,
	}

	if err = provider.Send(message); err != nil {
		errorHandler(w, http.StatusBadRequest, err, receiverConf.Provider)
		return
	}

	requestTotal.WithLabelValues("200", receiverConf.Provider).Inc()
}

func (h handlers) Reload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodPost {
		log.Println("Loading configuration file", *configFile)
		if err := LoadConfig(*configFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
	}
}
