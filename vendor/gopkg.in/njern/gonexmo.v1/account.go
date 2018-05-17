package nexmo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Account represents the user's account. Used when retrieving e.g current
// balance.
type Account struct {
	client *Client
}

// GetBalance retrieves the current balance of your Nexmo account in Euros (€)
func (nexmo *Account) GetBalance() (float64, error) {
	// Declare this locally, since we are only going to return a float64.
	type AccountBalance struct {
		Value float64 `json:"value"`
	}

	var accBalance *AccountBalance

	r, reqErr := http.NewRequest("GET", apiRoot+"/account/get-balance/"+
		nexmo.client.apiKey+"/"+nexmo.client.apiSecret, nil)

	if reqErr != nil {
		return 0.0, reqErr
	}

	r.Header.Add("Accept", "application/json")

	resp, err := nexmo.client.HttpClient.Do(r)
	if err != nil {
		return 0.0, err
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return 0.0, readErr
	}

	err = json.Unmarshal(body, &accBalance)
	if err != nil {
		return 0.0, err
	}

	return accBalance.Value, nil
}
