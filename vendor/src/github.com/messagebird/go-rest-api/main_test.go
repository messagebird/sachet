package messagebird

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var mbClient *Client
var mbServer *httptest.Server
var mbServerResponseCode int
var mbServerResponseBody []byte

var accessKeyErrorObject []byte = []byte(`{
  "errors":[
    {
      "code":2,
      "description":"Request not allowed (incorrect access_key)",
      "parameter":"access_key"
    }
  ]
}`)

// startFauxServer sets up a fake HTTPS server which the testing client will
// connect to, instead of the actual https://rest.messagebird.com URL.
func startFauxServer() {
	mbServer = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(mbServerResponseCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write(mbServerResponseBody)
	}))

	transport := &http.Transport{
		DialTLS: func(netw, addr string) (net.Conn, error) {
			conn, err := tls.Dial(netw, mbServer.Listener.Addr().String(), &tls.Config{InsecureSkipVerify: true})
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}

	mbClient = &Client{
		AccessKey:  "test_gshuPaZoeEG6ovbc8M79w0QyM",
		HTTPClient: &http.Client{Transport: transport},
	}

	if testing.Verbose() {
		mbClient.DebugLog = log.New(os.Stdout, "DEBUG", log.Lshortfile)
	}
}

// stopFauxServer is called after testing is done and stops the fake HTTPS
// server.
func stopFauxServer() {
	mbServer.Close()
}

func TestMain(m *testing.M) {
	startFauxServer()
	exitCode := m.Run()
	stopFauxServer()

	os.Exit(exitCode)
}

// SetServerResponse sets the response and HTTP status code the fake
// MessageBird server should return.
func SetServerResponse(statusCode int, response []byte) {
	mbServerResponseCode, mbServerResponseBody = statusCode, response
}
