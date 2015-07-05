package application

import (
	"testing"
)

func TestCentralDefaults(t *testing.T) {
	app := &Application{}
	app.Settings = make(map[string]string)
	app.Settings["MCRHUB_MODE"] = "central"

	if !app.IsCentralMode() {
		t.Error("By default, application should be in central mode.")
	}
}

func TestCentralTokens(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Fatalf("Failed to create application. Error: %v", err)
	}
	app.Settings["MCRHUB_TOKENS_DIR"] = "$GOPATH/src/github.com/didip/mcrouter-hub/tests/tokens"

	tokens := app.Tokens()
	if len(tokens) == 0 {
		t.Errorf("Tokens should not be empty.")
	}
	for _, token := range tokens {
		found := false
		for _, expectedToken := range []string{"aaa", "bbb", "ccc"} {
			if token == expectedToken {
				found = true
			}
		}
		if !found {
			t.Errorf("Unexpected token: %v", token)
		}
	}
}
