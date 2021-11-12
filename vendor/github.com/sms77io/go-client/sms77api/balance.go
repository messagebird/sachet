package sms77api

import "strconv"

type BalanceResource resource

func (api *BalanceResource) Get() (*float64, error) {
	res, err := api.client.request("balance", "GET", nil)

	if err != nil {
		return nil, err
	}

	float, err := strconv.ParseFloat(res, 64)
	if err != nil {
		return nil, err
	}

	return &float, nil
}
