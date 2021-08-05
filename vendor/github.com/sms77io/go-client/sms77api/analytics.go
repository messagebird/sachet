package sms77api

import (
	"encoding/json"
)

type AnalyticsParams struct {
	End         string           `json:"end,omitempty"`
	GroupBy     AnalyticsGroupBy `json:"group_by,omitempty"`
	Label       string           `json:"label,omitempty"`
	Start       string           `json:"start,omitempty"`
	Subaccounts string           `json:"subaccounts,omitempty"`
}

type Analytics struct {
	Account *string `json:"account"`
	Country *string `json:"country"`
	Date    *string `json:"date"`
	Label   *string `json:"label"`

	Direct   int     `json:"direct"`
	Economy  int     `json:"economy"`
	Hlr      int     `json:"hlr"`
	Inbound  int     `json:"inbound"`
	Mnp      int     `json:"mnp"`
	Voice    int     `json:"voice"`
	UsageEur float64 `json:"usage_eur"`
}

type AnalyticsResource resource

type AnalyticsGroupBy string

const (
	AnalyticsGroupByCountry    AnalyticsGroupBy = "country"
	AnalyticsGroupByDate       AnalyticsGroupBy = "date"
	AnalyticsGroupByLabel      AnalyticsGroupBy = "label"
	AnalyticsGroupBySubaccount AnalyticsGroupBy = "subaccount"
)

func (api *AnalyticsResource) Get(p *AnalyticsParams) (o []Analytics, err error) {
	res, err := api.client.request("analytics", "GET", p)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &o)

	if nil != err {
		return nil, err
	}

	return o, nil
}
