package messagebird

import (
	"testing"
	"time"
)

var hlrObject []byte = []byte(`{
  "id":"27978c50354a93ca0ca8de6h54340177",
  "href":"https:\/\/rest.messagebird.com\/hlr\/27978c50354a93ca0ca8de6h54340177",
  "msisdn":31612345678,
  "network":20406,
  "reference":"MyReference",
  "status":"sent",
  "createdDatetime":"2015-01-04T13:14:08+00:00",
  "statusDatetime":"2015-01-04T13:14:09+00:00"
}`)

func assertHLRObject(t *testing.T, hlr *HLR) {
	if hlr.Id != "27978c50354a93ca0ca8de6h54340177" {
		t.Errorf("Unexpected result for HLR Id: %s", hlr.Id)
	}

	if hlr.HRef != "https://rest.messagebird.com/hlr/27978c50354a93ca0ca8de6h54340177" {
		t.Errorf("Unexpected HLR href: %s", hlr.HRef)
	}

	if hlr.MSISDN != 31612345678 {
		t.Errorf("Unexpected HLR msisdn: %d", hlr.MSISDN)
	}

	if hlr.Network != 20406 {
		t.Errorf("Unexpected HLR network: %d", hlr.Network)
	}

	if hlr.Reference != "MyReference" {
		t.Errorf("Unexpected HLR reference: %s", hlr.Reference)
	}

	if hlr.Status != "sent" {
		t.Errorf("Unexpected HLR status: %s", hlr.Status)
	}

	if hlr.CreatedDatetime == nil || hlr.CreatedDatetime.Format(time.RFC3339) != "2015-01-04T13:14:08Z" {
		t.Errorf("Unexpected HLR created datetime: %s", hlr.CreatedDatetime.Format(time.RFC3339))
	}

	if hlr.StatusDatetime == nil || hlr.StatusDatetime.Format(time.RFC3339) != "2015-01-04T13:14:09Z" {
		t.Errorf("Unexpected HLR status datetime: %s", hlr.StatusDatetime.Format(time.RFC3339))
	}
}

func TestHLR(t *testing.T) {
	SetServerResponse(200, hlrObject)

	hlr, err := mbClient.HLR("27978c50354a93ca0ca8de6h54340177")
	if err != nil {
		t.Fatalf("Didn't expect an error while requesting a HLR: %s", err)
	}

	assertHLRObject(t, hlr)
}

func TestNewHLR(t *testing.T) {
	SetServerResponse(200, hlrObject)

	hlr, err := mbClient.NewHLR("31612345678", "MyReference")
	if err != nil {
		t.Fatalf("Didn't expect an error while creating a new HLR: %s", err)
	}

	assertHLRObject(t, hlr)
}

func TestHLRError(t *testing.T) {
	SetServerResponse(405, accessKeyErrorObject)

	hlr, err := mbClient.HLR("dummy_hlr_id")
	if err != ErrResponse {
		t.Fatalf("Expected ErrResponse to be returned, instead I got %s", err)
	}

	if len(hlr.Errors) != 1 {
		t.Fatalf("Unexpected number of errors: %d", len(hlr.Errors))
	}

	if hlr.Errors[0].Code != 2 {
		t.Errorf("Unexpected error code: %d", hlr.Errors[0].Code)
	}

	if hlr.Errors[0].Parameter != "access_key" {
		t.Errorf("Unexpected error parameter: %s", hlr.Errors[0].Parameter)
	}
}
