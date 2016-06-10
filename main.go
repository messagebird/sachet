package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/alertmanager/template"
)

var (
	listenAddress = flag.String("listen-address", ":9876", "The address to listen on for HTTP requests.")
)

func main() {

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		fmt.Println("new request")
		// content, _ := ioutil.ReadAll(r.Body)
		// fmt.Println(string(content))

		// https://godoc.org/github.com/prometheus/alertmanager/template#Data
		data := template.Data{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			fmt.Println(err.Error())
			return
		}

		receiverConf := receiverConfByReceiver(data.Receiver)
		provider := providerByName(receiverConf.Provider)

		for _, alert := range data.Alerts {

			text := receiverConf.Text
			if text == "" {
				tmpl(&alert)
			}

			message := Message{
				To:   receiverConf.To,
				From: receiverConf.From,
				Text: text,
			}
			provider.Send(message)
		}
	})

	http.Handle("/metrics", prometheus.Handler())

	if os.Getenv("PORT") != "" {
		*listenAddress = ":" + os.Getenv("PORT")
	}

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func tmpl(alert *template.Alert) string {
	// TODO "a message made of parts of alert"
	// https: //godoc.org/github.com/prometheus/alertmanager/template#Alert
	return alert.Status + ": " + alert.Labels["alertname"]
}

// receiverConfByReceiver loops the receiver conf list and returns the first instance with that name
func receiverConfByReceiver(name string) *ReceiverConf {
	for i, _ := range config.Receivers {
		rc := &config.Receivers[i]
		if rc.Name == name {
			return rc
		}
	}
	return nil
}

type Message struct {
	To   string
	From string
	Text string
}

type Provider interface {
	Send(message Message)
}

func providerByName(name string) Provider {
	switch name {
	case "messagebird":
		return &MessageBird{}
	case "nexmo":
		return &Nexmo{}
	}

	return nil
}
