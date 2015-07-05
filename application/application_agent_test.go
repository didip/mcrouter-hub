package application

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("MCROUTER_ADDR", "localhost:5000")
	os.Setenv("MCROUTER_CONFIG_FILE", os.ExpandEnv("$GOPATH/src/github.com/didip/mcrouter-hub/tests/mcrouter.json"))
	os.Setenv("MCRHUB_MODE", "agent")
}

func TestAgentDefaults(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Fatalf("Failed to create application. Error: %v", err)
	}

	if !app.IsAgentMode() {
		t.Errorf("By default, application should be in agent mode. Settings: %v", app.Settings)
	}
	if !app.IsReadOnly() {
		t.Error("By default, application should be in read only mode.")
	}
	if app.Settings["MCRHUB_REPORT_INTERVAL"] == "" {
		t.Error("By default, MCRHUB_REPORT_INTERVAL should not be empty.")
	}
	if app.Settings["MCRHUB_ADDR"] == "" {
		t.Error("By default, MCRHUB_ADDR should not be empty.")
	}
}

func TestAgentTokens(t *testing.T) {
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
