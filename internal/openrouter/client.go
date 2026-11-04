package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const requestTimeout = 2 * time.Minute

type Client struct {
	url    string
	apiKey string
	model  string
	http   *http.Client
}

func NewClient(url, apiKey, model string) *Client {
	return &Client{
		url:    url,
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{Timeout: requestTimeout},
	}
}

type Request struct {
	Instructions string
	Input        []json.RawMessage
	Tools        []map[string]any
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}

type Response struct {
	// Items holds the raw output items, which must be sent back unchanged as
	// input on the next request because the API keeps no conversation state.
	Items     []json.RawMessage
	Text      string
	ToolCalls []ToolCall
}

func UserMessage(text string) json.RawMessage {
	return marshalItem(map[string]any{
		"type": "message",
		"role": "user",
		"content": []map[string]string{
			{"type": "input_text", "text": text},
		},
	})
}

func UserImage(imageURL string) json.RawMessage {
	return marshalItem(map[string]any{
		"type": "message",
		"role": "user",
		"content": []map[string]string{
			{"type": "input_image", "image_url": imageURL, "detail": "auto"},
		},
	})
}

func ToolOutput(callID, output string) json.RawMessage {
	return marshalItem(map[string]string{
		"type":    "function_call_output",
		"call_id": callID,
		"output":  output,
	})
}

func marshalItem(item any) json.RawMessage {
	raw, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}
	return raw
}

func (c *Client) Respond(ctx context.Context, r Request) (Response, error) {
	body, err := json.Marshal(struct {
		Model        string            `json:"model"`
		Instructions string            `json:"instructions,omitempty"`
		Input        []json.RawMessage `json:"input"`
		Tools        []map[string]any  `json:"tools,omitempty"`
	}{c.model, r.Instructions, r.Input, r.Tools})
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("api returned %d: %s", resp.StatusCode, errorMessage(respBody))
	}
	return parseResponse(respBody)
}

type apiError struct {
	Message string `json:"message"`
}

func errorMessage(body []byte) string {
	var parsed struct {
		Error *apiError `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || parsed.Error == nil {
		return string(body)
	}
	return parsed.Error.Message
}

func parseResponse(body []byte) (Response, error) {
	var parsed struct {
		Error  *apiError         `json:"error"`
		Output []json.RawMessage `json:"output"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return Response{}, fmt.Errorf("decode response: %w", err)
	}
	if parsed.Error != nil {
		return Response{}, errors.New(parsed.Error.Message)
	}

	result := Response{Items: parsed.Output}
	var text strings.Builder
	for _, raw := range parsed.Output {
		var item struct {
			Type      string `json:"type"`
			CallID    string `json:"call_id"`
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
			Content   []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.Unmarshal(raw, &item); err != nil {
			return Response{}, fmt.Errorf("decode output item: %w", err)
		}

		switch item.Type {
		case "function_call":
			result.ToolCalls = append(result.ToolCalls, ToolCall{item.CallID, item.Name, item.Arguments})
		case "message":
			for _, part := range item.Content {
				if part.Type == "output_text" {
					text.WriteString(part.Text)
				}
			}
		}
	}
	result.Text = text.String()
	return result, nil
}
