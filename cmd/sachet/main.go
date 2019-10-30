package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/messagebird/sachet"
	"github.com/messagebird/sachet/provider/aspsms"
	"github.com/messagebird/sachet/provider/cm"
	"github.com/messagebird/sachet/provider/exotel"
	"github.com/messagebird/sachet/provider/freemobile"
	"github.com/messagebird/sachet/provider/infobip"
	"github.com/messagebird/sachet/provider/mediaburst"
	"github.com/messagebird/sachet/provider/messagebird"
	"github.com/messagebird/sachet/provider/nexmo"
	"github.com/messagebird/sachet/provider/otc"
	"github.com/messagebird/sachet/provider/pushbullet"
	"github.com/messagebird/sachet/provider/sipgate"
	"github.com/messagebird/sachet/provider/telegram"
	"github.com/messagebird/sachet/provider/turbosms"
	"github.com/messagebird/sachet/provider/twilio"

	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/heptiolabs/healthcheck"
)

var (
	listenAddress = flag.String("listen-address", ":9876", "The address to listen on for HTTP requests.")
	configFile    = flag.String("config", "config.yaml", "The configuration file")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := LoadConfig(*configFile); err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	http.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
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
			if len(data.Alerts) > 1 {
				labelAlerts := map[string]template.Alerts{
					"Firing":   data.Alerts.Firing(),
					"Resolved": data.Alerts.Resolved(),
				}
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
			} else if len(data.Alerts) == 1 {
				alert := data.Alerts[0]
				tuples := []string{}
				for k, v := range alert.Labels {
					tuples = append(tuples, k+"= "+v)
				}
				text = strings.ToUpper(data.Status) + " \n" + strings.Join(tuples, "\n")
			} else {
				text = "Alert \n" + strings.Join(data.CommonLabels.Values(), " | ")
			}
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
	})

	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method == "POST" {
			log.Println("Loading configuration file", *configFile)
			if err := LoadConfig(*configFile); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		}
	})

	hc := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "sachet")

	http.HandleFunc("/-/live", hc.LiveEndpoint)
	http.HandleFunc("/-/ready", hc.ReadyEndpoint)

	if os.Getenv("PORT") != "" {
		*listenAddress = ":" + os.Getenv("PORT")
	}

	log.Printf("Listening on %s", *listenAddress)

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

// receiverConfByReceiver loops the receiver conf list and returns the first instance with that name
func receiverConfByReceiver(name string) *ReceiverConf {
	for i := range config.Receivers {
		rc := &config.Receivers[i]
		if rc.Name == name {
			return rc
		}
	}
	return nil
}

func providerByName(name string) (sachet.Provider, error) {
	switch name {
	case "messagebird":
		return messagebird.NewMessageBird(config.Providers.MessageBird), nil
	case "nexmo":
		return nexmo.NewNexmo(config.Providers.Nexmo)
	case "twilio":
		return twilio.NewTwilio(config.Providers.Twilio), nil
	case "infobip":
		return infobip.NewInfobip(config.Providers.Infobip), nil
	case "turbosms":
		return turbosms.NewTurbosms(config.Providers.Turbosms), nil
	case "exotel":
		return exotel.NewExotel(config.Providers.Exotel), nil
	case "cm":
		return cm.NewCM(config.Providers.CM), nil
	case "telegram":
		return telegram.NewTelegram(config.Providers.Telegram)
	case "otc":
		return otc.NewOTC(config.Providers.OTC), nil
	case "mediaburst":
		return mediaburst.NewMediaBurst(config.Providers.MediaBurst), nil
	case "freemobile":
		return freemobile.NewFreeMobile(config.Providers.FreeMobile), nil
	case "aspsms":
		return aspsms.NewAspSms(config.Providers.AspSms), nil
	case "sipgate":
		return sipgate.NewSipgate(config.Providers.Sipgate), nil
	case "pushbullet":
		return pushbullet.NewPushbullet(config.Providers.Pushbullet), nil
	}

	return nil, fmt.Errorf("%s: Unknown provider", name)
}

func errorHandler(w http.ResponseWriter, status int, err error, provider string) {
	w.WriteHeader(status)

	data := struct {
		Error   bool
		Status  int
		Message string
	}{
		true,
		status,
		err.Error(),
	}
	// respond json
	bytes, _ := json.Marshal(data)
	json := string(bytes[:])
	fmt.Fprint(w, json)

	log.Println("Error: " + json)
	requestTotal.WithLabelValues(strconv.FormatInt(int64(status), 10), provider).Inc()
}
