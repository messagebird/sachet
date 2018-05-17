package twiliogo

import (
	"os"
	"testing"
)

var API_KEY string = os.Getenv("API_KEY")
var API_TOKEN string = os.Getenv("API_TOKEN")
var FROM_NUMBER string = os.Getenv("FROM_NUMBER")

var TEST_KEY string = os.Getenv("TEST_KEY")
var TEST_TOKEN string = os.Getenv("TEST_TOKEN")
var TEST_FROM_NUMBER string = os.Getenv("TEST_FROM_NUMBER")

var TO_NUMBER string = os.Getenv("TO_NUMBER")

func CheckTestEnv(t *testing.T) {
	if API_KEY == "" || API_TOKEN == "" ||
		TEST_KEY == "" || TEST_TOKEN == "" ||
		TEST_FROM_NUMBER == "" || FROM_NUMBER == "" ||
		TO_NUMBER == "" {
		t.SkipNow()
	}
}

func TestBuildUri(t *testing.T) {
	c := NewClient("abc", "")
	uri := c.buildUri("qzx")
	if uri != ROOT+"/"+VERSION+"/Accounts/abc/qzx" {
		t.Errorf("buildUri failed: got %s", uri)
	}
}
