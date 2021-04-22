package twiliogo

import (
	"encoding/json"
	"net/url"
)

type CallList struct {
	Client          Client
	Start           int    `json:"start"`
	Total           int    `json:"total"`
	NumPages        int    `json:"num_pages"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
	End             int    `json:"end"`
	Uri             string `json:"uri"`
	FirstPageUri    string `json:"first_page_uri"`
	LastPageUri     string `json:"last_page_uri"`
	NextPageUri     string `json:"next_page_uri"`
	PreviousPageUri string `json"previous_page_uri"`
	Calls           []Call `json:"calls"`
}

func GetCallList(client Client, optionals ...Optional) (*CallList, error) {
	var callList *CallList

	params := url.Values{}

	for _, optional := range optionals {
		param, value := optional.GetParam()
		params.Set(param, value)
	}

	body, err := client.get(nil, "/Calls.json")

	if err != nil {
		return nil, err
	}

	callList = new(CallList)
	callList.Client = client
	err = json.Unmarshal(body, callList)

	return callList, err
}

func (callList *CallList) GetCalls() []Call {
	return callList.Calls
}

func (currentCallList *CallList) HasNextPage() bool {
	return currentCallList.NextPageUri != ""
}

func (currentCallList *CallList) NextPage() (*CallList, error) {
	if !currentCallList.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return currentCallList.getPage(currentCallList.NextPageUri)
}

func (currentCallList *CallList) HasPreviousPage() bool {
	return currentCallList.PreviousPageUri != ""
}

func (currentCallList *CallList) PreviousPage() (*CallList, error) {
	if !currentCallList.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return currentCallList.getPage(currentCallList.NextPageUri)
}

func (currentCallList *CallList) FirstPage() (*CallList, error) {
	return currentCallList.getPage(currentCallList.FirstPageUri)
}

func (currentCallList *CallList) LastPage() (*CallList, error) {
	return currentCallList.getPage(currentCallList.LastPageUri)
}

func (currentCallList *CallList) getPage(uri string) (*CallList, error) {
	var callList *CallList

	client := currentCallList.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return callList, err
	}

	callList = new(CallList)
	callList.Client = client
	err = json.Unmarshal(body, callList)

	return callList, err
}
