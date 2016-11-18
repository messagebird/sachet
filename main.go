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

	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddress = flag.String("listen-address", ":9876", "The address to listen on for HTTP requests.")
	configFile    = flag.String("config", "config.yaml", "The configuration file")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	LoadConfig(*configFile)

	http.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// content, _ := ioutil.ReadAll(r.Body)
		// fmt.Println(string(content))

		// https://godoc.org/github.com/prometheus/alertmanager/template#Data
		data := template.Data{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			errorHandler(w, http.StatusBadRequest, err, "?")
			return
		}

		receiverConf := receiverConfByReceiver(data.Receiver)
		if receiverConf == nil {
			errorHandler(w, http.StatusBadRequest, fmt.Errorf("Receiver missing"), "?")
			return
		}
		provider := providerByName(receiverConf.Provider)
		if provider == nil {
			errorHandler(w, http.StatusBadRequest, fmt.Errorf("Cannot find provider implementation for '%s'", receiverConf.Provider), receiverConf.Provider)
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

		err := provider.Send(message)
		if err != nil {
			errorHandler(w, http.StatusBadRequest, err, receiverConf.Provider)
			return
		}

		requestTotal.WithLabelValues("200", receiverConf.Provider).Inc()
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
	Send(message Message) error
}

func providerByName(name string) Provider {
	switch name {
	case "messagebird":
		return &MessageBird{}
	case "nexmo":
		return &Nexmo{}
	case "twilio":
		return &Twilio{}
	}

	return nil
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
