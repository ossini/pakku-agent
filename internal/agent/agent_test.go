package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pakku-agent/internal/openrouter"
	"pakku-agent/internal/tools"
)

const (
	toolCallOutput = `{"output":[{"type":"function_call","call_id":"c1","name":"echo","arguments":"{}"}]}`
	textOutput     = `{"output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"fertig"}]}]}`
)

type recorder struct{ events []string }

func (r *recorder) Thinking() {}

func (r *recorder) ToolStarted(title, _ string) {
	r.events = append(r.events, "start:"+title)
}

func (r *recorder) ToolFinished(output string, err error) {
	if err != nil {
		output = err.Error()
	}
	r.events = append(r.events, "end:"+output)
}

type requestBody struct {
	Input []json.RawMessage `json:"input"`
}

func newAgent(t *testing.T, run func(context.Context, string) (tools.Result, error), responses ...string) (*Agent, *recorder, *[]requestBody) {
	t.Helper()

	requests := &[]requestBody{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body requestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
		}
		*requests = append(*requests, body)
		fmt.Fprint(w, responses[min(len(*requests), len(responses))-1])
	}))
	t.Cleanup(server.Close)

	manager := &tools.Manager{Tools: []tools.Tool{{Name: "echo", Title: "Echo", Run: run}}}
	observer := &recorder{}
	client := openrouter.NewClient(server.URL, "key", "model")
	return New(client, manager, 3, observer), observer, requests
}

func TestSendRunsToolAndReturnsFinalText(t *testing.T) {
	run := func(context.Context, string) (tools.Result, error) { return tools.Result{Text: "tool-result"}, nil }
	pakku, observer, requests := newAgent(t, run, toolCallOutput, textOutput)

	reply, err := pakku.Send(context.Background(), "hallo")
	if err != nil {
		t.Fatal(err)
	}
	if reply != "fertig" {
		t.Errorf("reply = %q", reply)
	}
	if got := strings.Join(observer.events, ","); got != "start:Echo,end:tool-result" {
		t.Errorf("events = %q", got)
	}

	input := (*requests)[1].Input
	if len(input) != 3 {
		t.Fatalf("second request has %d input items, want 3", len(input))
	}
	if !strings.Contains(string(input[2]), `"output":"tool-result"`) {
		t.Errorf("tool output item = %s", input[2])
	}
}

func TestSendAttachesImageAfterToolOutput(t *testing.T) {
	run := func(context.Context, string) (tools.Result, error) {
		return tools.Result{Text: "foto", ImageURL: "data:image/jpeg;base64,AAAA"}, nil
	}
	pakku, _, requests := newAgent(t, run, toolCallOutput, textOutput)

	if _, err := pakku.Send(context.Background(), "hallo"); err != nil {
		t.Fatal(err)
	}
	input := (*requests)[1].Input
	if len(input) != 4 {
		t.Fatalf("second request has %d input items, want 4", len(input))
	}
	if !strings.Contains(string(input[2]), "function_call_output") {
		t.Errorf("tool output item = %s", input[2])
	}
	if !strings.Contains(string(input[3]), `"image_url":"data:image/jpeg;base64,AAAA"`) {
		t.Errorf("image item = %s", input[3])
	}
}

func TestSendReportsToolErrorToModel(t *testing.T) {
	run := func(context.Context, string) (tools.Result, error) { return tools.Result{}, errors.New("kaputt") }
	pakku, _, requests := newAgent(t, run, toolCallOutput, textOutput)

	if _, err := pakku.Send(context.Background(), "hallo"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string((*requests)[1].Input[2]), "error: kaputt") {
		t.Errorf("tool output item = %s", (*requests)[1].Input[2])
	}
}

func TestSendKeepsHistoryAcrossTurns(t *testing.T) {
	pakku, _, requests := newAgent(t, nil, textOutput)

	for range 2 {
		if _, err := pakku.Send(context.Background(), "hallo"); err != nil {
			t.Fatal(err)
		}
	}
	if got := len((*requests)[1].Input); got != 3 {
		t.Errorf("second turn has %d input items, want 3", got)
	}
}

func TestFailedTurnLeavesHistoryUntouched(t *testing.T) {
	run := func(context.Context, string) (tools.Result, error) { return tools.Result{Text: "ok"}, nil }
	pakku, _, requests := newAgent(t, run, toolCallOutput)

	if _, err := pakku.Send(context.Background(), "hallo"); err == nil {
		t.Fatal("expected step limit error")
	}
	if len(*requests) != 3 {
		t.Errorf("made %d requests, want 3", len(*requests))
	}
	if len(pakku.history) != 0 {
		t.Errorf("history has %d items, want 0", len(pakku.history))
	}
}
