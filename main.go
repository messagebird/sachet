package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/alertmanager/template"
)

var (
	listenAddress = flag.String("listen-address", ":7733", "The address to listen on for HTTP requests.")
)

func main() {

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		fmt.Println("new request")
		content, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(content))

		data := template.Data{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&data); err != nil {
			fmt.Println(err.Error())
			return
		}

	})

	http.Handle("/metrics", prometheus.Handler())

	if os.Getenv("PORT") != "" {
		*listenAddress = ":" + os.Getenv("PORT")
	}

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
