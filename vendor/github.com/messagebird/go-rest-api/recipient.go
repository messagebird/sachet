package messagebird

import "time"

type Recipient struct {
	Recipient      int
	Status         string
	StatusDatetime *time.Time
}

type Recipients struct {
	TotalCount               int
	TotalSentCount           int
	TotalDeliveredCount      int
	TotalDeliveryFailedCount int
	Items                    []Recipient
}
