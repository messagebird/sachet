package messagebird

import "time"

// Recipient struct holds information for a single msisdn with status details.
type Recipient struct {
	Recipient      int
	Status         string
	StatusDatetime *time.Time
}

// Recipients holds a collection of Recepient structs along with send stats.
type Recipients struct {
	TotalCount               int
	TotalSentCount           int
	TotalDeliveredCount      int
	TotalDeliveryFailedCount int
	Items                    []Recipient
}
