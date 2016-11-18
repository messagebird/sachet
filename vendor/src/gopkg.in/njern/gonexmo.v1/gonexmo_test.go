package nexmo

import (
	"fmt"
	"os"
	"testing"
)

var (
	API_KEY, API_SECRET, TEST_PHONE_NUMBER, TEST_FROM string
)

func init() {
	API_KEY = os.Getenv("NEXMO_KEY")
	if API_KEY == "" {
		fmt.Println("No API key specified. Please set NEXMO_KEY")
		os.Exit(1)
	}

	API_SECRET = os.Getenv("NEXMO_SECRET")
	if API_SECRET == "" {
		fmt.Println("No API secret specified. Please set NEXMO_SECRET")
		os.Exit(1)
	}

	TEST_PHONE_NUMBER = os.Getenv("NEXMO_NUM")

	// Set a custom from value, or use the default. If you get error 15 when
	// sending a message ("Illegal Sender Address - rejected") try setting this
	// to your nexmo phone number.
	TEST_FROM = os.Getenv("NEXMO_FROM")
	if TEST_FROM == "" {
		TEST_FROM = "gonexmo"
	}
}

func TestNexmoCreation(t *testing.T) {
	_, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}
}
