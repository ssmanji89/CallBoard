package main

import "testing"

func TestTailnetOf(t *testing.T) {
	up := `{"BackendState":"Running","Self":{"DNSName":"mac.tail1.ts.net.","UserID":42},"User":{"42":{"LoginName":"me@example.com"},"7":{"LoginName":"other@example.com"}}}`
	got, err := tailnetOf([]byte(up))
	if err != nil || got.Host != "mac.tail1.ts.net" || got.Login != "me@example.com" {
		t.Errorf("tailnetOf: %+v, %v", got, err)
	}
	for _, bad := range []string{
		`{"BackendState":"Stopped","Self":{"DNSName":"mac.tail1.ts.net.","UserID":42},"User":{"42":{"LoginName":"me@example.com"}}}`,
		`{"BackendState":"Running","Self":{"DNSName":"mac.tail1.ts.net.","UserID":42},"User":{}}`,
		`not json`,
	} {
		if _, err := tailnetOf([]byte(bad)); err == nil {
			t.Errorf("tailnetOf(%s): want an error", bad)
		}
	}
}
