package twiliogo

import (
	"net/url"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (client *MockClient) AccountSid() string {
	return "AC3FakeClient"
}

func (client *MockClient) AuthToken() string {
	return "98h4hfaketoken"
}

func (client *MockClient) RootUrl() string {
	return "http://test.com/fake"
}

func (client *MockClient) get(params url.Values, uri string) ([]byte, error) {
	args := client.Mock.Called(params, uri)
	return args.Get(0).([]byte), args.Error(1)
}

func (client *MockClient) post(params url.Values, uri string) ([]byte, error) {
	args := client.Mock.Called(params, uri)
	return args.Get(0).([]byte), args.Error(1)
}

func (client *MockClient) delete(uri string) error {
	args := client.Mock.Called(nil, uri)
	return args.Error(1)
}
