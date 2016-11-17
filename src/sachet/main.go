package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddress = flag.String("listen-address", ":9876", "The address to listen on for HTTP requests.")
	configFile    = flag.String("config", "config.yaml", "The configuration file")
)

func main() {
	flag.Parse()

	LoadConfig(*configFile)

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
		if receiverConf == nil {
			// fmt.Println("no receiver")
			return
		}
		provider := providerByName(receiverConf.Provider)
		if provider == nil {
			// fmt.Println("no provider")
			return
		}

		// Concatenate common labels to form the alert string.
		text := strings.Join(data.CommonLabels.Values(), " | ")
		if len(text) > 160 {
			text = text[:160]
		}

		message := Message{
			To:   receiverConf.To,
			From: receiverConf.From,
			Text: text,
		}
		provider.Send(message)
	})

	http.Handle("/metrics", prometheus.Handler())

	if os.Getenv("PORT") != "" {
		*listenAddress = ":" + os.Getenv("PORT")
	}

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

type Message struct {
	To   []string
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
