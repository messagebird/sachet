package messagebird

import "testing"

var balanceObject []byte = []byte(`{
  "payment":"prepaid",
  "type":"credits",
  "amount":9.2
}`)

const Epsilon float32 = 0.001

func cmpFloat32(a, b float32) bool {
	return (a-b) < Epsilon && (b-a) < Epsilon
}

func TestBalance(t *testing.T) {
	SetServerResponse(200, balanceObject)

	balance, err := mbClient.Balance()
	if err != nil {
		t.Fatalf("Didn't expect error while fetching the balance: %s", err)
	}

	if balance.Payment != "prepaid" {
		t.Errorf("Unexpected balance payment: %s", balance.Payment)
	}
	if balance.Type != "credits" {
		t.Errorf("Unexpected balance type: %s", balance.Type)
	}
	if !cmpFloat32(balance.Amount, 9.2) {
		t.Errorf("Unexpected balance amount: %.2f", balance.Amount)
	}
}

func TestBalanceError(t *testing.T) {
	SetServerResponse(405, accessKeyErrorObject)

	balance, err := mbClient.Balance()
	if err != ErrResponse {
		t.Fatalf("Expected ErrResponse to be returned, instead I got %s", err)
	}

	if len(balance.Errors) != 1 {
		t.Fatalf("Unexpected number of errors: %d", len(balance.Errors))
	}

	if balance.Errors[0].Code != 2 {
		t.Errorf("Unexpected error code: %d", balance.Errors[0].Code)
	}

	if balance.Errors[0].Parameter != "access_key" {
		t.Errorf("Unexpected error parameter: %s", balance.Errors[0].Parameter)
	}
}
