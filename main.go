package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddress = flag.String("listen-address", ":7733", "The address to listen on for HTTP requests.")
)

func main() {

	flag.Parse()

	http.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.Handle("/metrics", prometheus.Handler())

	if os.Getenv("PORT") != "" {
		*listenAddress = ":" + os.Getenv("PORT")
	}

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
