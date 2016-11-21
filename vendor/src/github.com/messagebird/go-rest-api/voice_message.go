package messagebird

import (
	"net/url"
	"strconv"
	"time"
)

type VoiceMessage struct {
	Id                string
	HRef              string
	Originator        string
	Body              string
	Reference         string
	Language          string
	Voice             string
	Repeat            int
	IfMachine         string
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        Recipients
	Errors            []Error
}

type VoiceMessageParams struct {
	Originator        string
	Reference         string
	Language          string
	Voice             string
	Repeat            int
	IfMachine         string
	ScheduledDatetime time.Time
}

// paramsForVoiceMessage converts the specified VoiceMessageParams struct to a
// url.Values pointer and returns it.
func paramsForVoiceMessage(params *VoiceMessageParams) *url.Values {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams
	}

	if params.Originator != "" {
		urlParams.Set("originator", params.Originator)
	}
	if params.Reference != "" {
		urlParams.Set("reference", params.Reference)
	}
	if params.Language != "" {
		urlParams.Set("language", params.Language)
	}
	if params.Voice != "" {
		urlParams.Set("voice", params.Voice)
	}

	// A repeat value of 1 actually means "play it once", not "repeat it once"
	// So only set the repeat value when it's larger than 1.
	if params.Repeat > 1 {
		urlParams.Set("repeat", strconv.Itoa(params.Repeat))
	}
	if params.IfMachine != "" {
		urlParams.Set("ifMachine", params.IfMachine)
	}
	if params.ScheduledDatetime.Unix() > 0 {
		urlParams.Set("scheduledDatetime", params.ScheduledDatetime.Format(time.RFC3339))
	}

	return urlParams
}
