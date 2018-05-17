package nexmo

import "testing"

var isTrustedIPTests = []struct {
	ip   string
	want bool
}{
	{"174.37.245.33", true},
	{"174.36.197.193", true},
	{"173.193.199.17", true},
	{"119.81.44.1", true},
	{"50.0.2.241", false},
	{"17.178.96.59", false},
}

func TestIsTrustedIP(t *testing.T) {
	for _, test := range isTrustedIPTests {
		got := IsTrustedIP(test.ip)
		if got != test.want {
			t.Errorf("IsTrustedIP(%s) = %v, want %v",
				test.ip, got, test.want)
		}
	}
}
