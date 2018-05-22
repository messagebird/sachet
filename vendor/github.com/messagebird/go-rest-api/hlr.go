package messagebird

import (
	"errors"
	"time"
)

// HLR stands for Home Location Register.
// Contains information about the subscribers identity, telephone number, the associated services and general information about the location of the subscriber
type HLR struct {
	ID              string
	HRef            string
	MSISDN          int
	Network         int
	Reference       string
	Status          string
	Details         map[string]interface{}
	CreatedDatetime *time.Time
	StatusDatetime  *time.Time
	Errors          []Error
}

// HLRList represents a list of HLR requests.
type HLRList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      map[string]*string
	Items      []HLR
}

type hlrRequest struct {
	MSISDN    string `json:"msisdn"`
	Reference string `json:"reference"`
}

func requestDataForHLR(msisdn string, reference string) (*hlrRequest, error) {
	if msisdn == "" {
		return nil, errors.New("msisdn is required")
	}
	if reference == "" {
		return nil, errors.New("reference is required")
	}

	request := &hlrRequest{
		MSISDN:    msisdn,
		Reference: reference,
	}

	return request, nil
}
