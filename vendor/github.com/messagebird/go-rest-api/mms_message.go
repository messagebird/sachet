package messagebird

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

// MMSMessage represents a MMS Message.
type MMSMessage struct {
	ID                string
	HRef              string
	Direction         string
	Originator        string
	Body              string
	Reference         string
	Subject           string
	MediaUrls         []string
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        Recipients
	Errors            []Error
}

// MMSMessageParams represents the parameters that can be supplied when creating
// a request.
type MMSMessageParams struct {
	Body              string
	MediaUrls         []string
	Subject           string
	Reference         string
	ScheduledDatetime time.Time
}

// paramsForMMSMessage converts the specified MMSMessageParams struct to a
// url.Values pointer and returns it.
func paramsForMMSMessage(params *MMSMessageParams) (*url.Values, error) {
	urlParams := &url.Values{}

	if params.Body == "" && params.MediaUrls == nil {
		return nil, errors.New("Body or MediaUrls is required")
	}
	if params.Body != "" {
		urlParams.Set("body", params.Body)
	}
	if params.MediaUrls != nil {
		urlParams.Set("mediaUrls[]", strings.Join(params.MediaUrls, ","))
	}
	if params.Subject != "" {
		urlParams.Set("subject", params.Subject)
	}
	if params.Reference != "" {
		urlParams.Set("reference", params.Reference)
	}
	if params.ScheduledDatetime.Unix() > 0 {
		urlParams.Set("scheduledDatetime", params.ScheduledDatetime.Format(time.RFC3339))
	}

	return urlParams, nil
}
