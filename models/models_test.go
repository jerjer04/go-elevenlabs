package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

func TestList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/models" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"model_id": "eleven_multilingual_v2",
				"name": "Eleven Multilingual v2",
				"can_be_finetuned": false,
				"can_do_text_to_speech": true,
				"can_do_voice_conversion": false,
				"can_use_style": true,
				"can_use_speaker_boost": true,
				"serves_pro_voices": false,
				"token_cost_factor": 1.0
			}
		]`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(result.Models) != 1 {
		t.Errorf("expected 1 model, got %d", len(result.Models))
	}
	if result.Models[0].ModelID != "eleven_multilingual_v2" {
		t.Errorf("unexpected model ID: %s", result.Models[0].ModelID)
	}
	if !result.Models[0].CanDoTextToSpeech {
		t.Error("expected model to support text-to-speech")
	}
}
