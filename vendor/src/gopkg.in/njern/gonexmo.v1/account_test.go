package nexmo

import (
	"testing"
)

func TestGetAccountBalance(t *testing.T) {
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	balance, err := nexmo.Account.GetBalance()
	if err != nil {
		t.Error("Failed to get account balance with error:", err)
	}

	t.Log("Got account balance: ", balance, "€")
}
