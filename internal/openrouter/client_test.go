package openrouter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseResponse(t *testing.T) {
	body := `{"output":[
		{"type":"reasoning","id":"r1"},
		{"type":"function_call","call_id":"c1","name":"get_time","arguments":"{}"},
		{"type":"message","content":[{"type":"output_text","text":"Hallo "},{"type":"output_text","text":"Welt"}]}
	]}`

	resp, err := parseResponse([]byte(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != "Hallo Welt" {
		t.Errorf("Text = %q", resp.Text)
	}
	if len(resp.Items) != 3 {
		t.Errorf("Items = %d, want 3", len(resp.Items))
	}
	if len(resp.ToolCalls) != 1 || resp.ToolCalls[0] != (ToolCall{"c1", "get_time", "{}"}) {
		t.Errorf("ToolCalls = %+v", resp.ToolCalls)
	}
}

func TestParseResponseReturnsEmbeddedError(t *testing.T) {
	if _, err := parseResponse([]byte(`{"error":{"message":"overloaded"}}`)); err == nil || err.Error() != "overloaded" {
		t.Errorf("err = %v", err)
	}
}

func TestRespondReportsStatusAndMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer key" {
			t.Errorf("Authorization = %q", got)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"invalid key"}}`))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "key", "model").Respond(context.Background(), Request{})
	if err == nil || !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "invalid key") {
		t.Errorf("err = %v", err)
	}
}
