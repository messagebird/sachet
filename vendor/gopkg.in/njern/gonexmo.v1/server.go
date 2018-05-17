package nexmo

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MessageType int

const (
	TextMessage = iota + 1
	UnicodeMessage
	BinaryMessage
)

var messageTypeMap = map[string]MessageType{
	"text":    TextMessage,
	"unicdoe": UnicodeMessage,
	"binary":  BinaryMessage,
}

var messageTypeIntMap = map[MessageType]string{
	TextMessage:    "text",
	UnicodeMessage: "unicode",
	BinaryMessage:  "binary",
}

func (m MessageType) String() string {
	if m < 1 || m > 3 {
		return "undefined"
	}
	return messageTypeIntMap[m]
}

// RecvdMessage represents a message that was received from the Nexmo API.
type RecvdMessage struct {
	// Expected values are "text" or "binary".
	Type MessageType

	// Recipient number (your long virtual number).
	To string

	// Sender ID.
	MSISDN string

	// Optional unique identifier of a mobile network MCCMNC.
	NetworkCode string

	// Nexmo message ID.
	ID string

	// Time when Nexmo started to push the message to you.
	Timestamp time.Time

	// Parameters for conactenated messages:
	Concatenated bool // Set to true if a MO concatenated message is detected.
	Concat       struct {

		// Transaction reference. All message parts will share the same
		//transaction reference.
		Reference string

		// Total number of parts in this concatenated message set.
		Total int

		// The part number of this message withing the set (starts at 1).
		Part int
	}

	// When Type == text:
	Text string // Content of the message

	Keyword string // First word in the message body, typically used with short codes
	// When type == binary:

	// Content of the message.
	Data []byte

	// User Data Header.
	UDH []byte
}

// DeliveryReceipt is a delivery receipt for a single SMS sent via the Nexmo API
type DeliveryReceipt struct {
	To              string    `json:"to"`
	NetworkCode     string    `json:"network-code"`
	MessageID       string    `json:"messageId"`
	MSISDN          string    `json:"msisdn"`
	Status          string    `json:"status"`
	ErrorCode       string    `json:"err-code"`
	Price           string    `json:"price"`
	SCTS            time.Time `json:"scts"`
	Timestamp       time.Time `json:"message-timestamp"`
	ClientReference string    `json:"client-ref"`
}

// NewDeliveryHandler creates a new http.HandlerFunc that can be used to listen
// for deivery receipts from the Nexmo server. Any receipts received will be
// decoded nad passed to the out chan.
func NewDeliveryHandler(out chan *DeliveryReceipt, verifyIPs bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if verifyIPs {
			// Check if the request came from Nexmo
			host, _, err := net.SplitHostPort(req.RemoteAddr)
			if !IsTrustedIP(host) || err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		var err error
		// Check if the query is empty. If it is, it's just Nexmo
		// making sure our service is up, so we don't want to return
		// an error.
		if req.URL.RawQuery == "" {
			return
		}

		req.ParseForm()
		// Decode the form data
		m := new(DeliveryReceipt)

		m.To = req.FormValue("to")
		m.NetworkCode = req.FormValue("network-code")
		m.MessageID = req.FormValue("messageId")
		m.MSISDN = req.FormValue("msisdn")
		m.Status = req.FormValue("status")
		m.ErrorCode = req.FormValue("err-code")
		m.Price = req.FormValue("price")
		m.ClientReference = req.FormValue("client-ref")

		t, err := url.QueryUnescape(req.FormValue("scts"))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Convert the timestamp to a time.Time.
		timestamp, err := time.Parse("0601021504", t)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.SCTS = timestamp

		t, err = url.QueryUnescape(req.FormValue("message-timestamp"))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Convert the timestamp to a time.Time.
		timestamp, err = time.Parse("2006-01-02 15:04:05", t)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.Timestamp = timestamp

		// Pass it out on the chan
		out <- m
	}

}

// NewMessageHandler creates a new http.HandlerFunc that can be used to listen
// for new messages from the Nexmo server. Any new messages received will be
// decoded and passed to the out chan.
func NewMessageHandler(out chan *RecvdMessage, verifyIPs bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if verifyIPs {
			// Check if the request came from Nexmo
			host, _, err := net.SplitHostPort(req.RemoteAddr)
			if !IsTrustedIP(host) || err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		var err error

		// Check if the query is empty. If it is, it's just Nexmo
		// making sure our service is up, so we don't want to return
		// an error.
		if req.URL.RawQuery == "" {
			return
		}

		req.ParseForm()
		// Decode the form data
		m := new(RecvdMessage)
		switch req.FormValue("type") {
		case "text":
			m.Text, err = url.QueryUnescape(req.FormValue("text"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Type = TextMessage
		case "unicode":
			m.Text, err = url.QueryUnescape(req.FormValue("text"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Type = UnicodeMessage

			// TODO: I have no idea if this data stuff works, as I'm unable to
			// send data SMS messages.
		case "binary":
			data, err := url.QueryUnescape(req.FormValue("data"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Data = []byte(data)

			udh, err := url.QueryUnescape(req.FormValue("udh"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.UDH = []byte(udh)
			m.Type = BinaryMessage

		default:
			//error
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.To = req.FormValue("to")
		m.MSISDN = req.FormValue("msisdn")
		m.NetworkCode = req.FormValue("network-code")
		m.ID = req.FormValue("messageId")

		m.Keyword = req.FormValue("keyword")
		t, err := url.QueryUnescape(req.FormValue("message-timestamp"))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Convert the timestamp to a time.Time.
		timestamp, err := time.Parse("2006-01-02 15:04:05", t)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.Timestamp = timestamp

		// TODO: I don't know if this works as I've been unable to send an SMS
		// message longer than 160 characters that doesn't get concatenated
		// automatically.
		if req.FormValue("concat") == "true" {
			m.Concatenated = true
			m.Concat.Reference = req.FormValue("concat-ref")
			m.Concat.Total, err = strconv.Atoi(req.FormValue("concat-total"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Concat.Part, err = strconv.Atoi(req.FormValue("concat-part"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		// Pass it out on the chan
		out <- m
	}

}
