package tts

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

func TestConvert(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/v1/text-to-speech/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("xi-api-key") == "" {
			t.Error("missing xi-api-key header")
		}

		// Return mock audio data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock audio data"))
	}))
	defer server.Close()

	// Create client with test server
	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	// Test request
	req := ConvertRequest{
		Text:    "Hello world",
		ModelID: "eleven_flash_v2_5",
		VoiceSettings: &VoiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0.75,
		},
	}

	opts := &ConvertOptions{
		OptimizeStreamingLatency: 0,
		OutputFormat:             "mp3_44100_128",
	}

	result, err := service.Convert(context.Background(), "test-voice-id", req, opts)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}
	defer result.Close()

	// Read response
	data, err := io.ReadAll(result)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if string(data) != "mock audio data" {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestConvertWithTimestamps(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/with-timestamps") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"audio_base64": "bW9jayBhdWRpbyBkYXRh",
			"alignment": {
				"characters": ["H", "e", "l", "l", "o"],
				"character_start_times_seconds": [0.0, 0.1, 0.2, 0.3, 0.4],
				"character_end_times_seconds": [0.1, 0.2, 0.3, 0.4, 0.5]
			}
		}`))
	}))
	defer server.Close()

	// Create client with test server
	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	// Test request
	req := ConvertRequest{
		Text: "Hello",
	}

	result, err := service.ConvertWithTimestamps(context.Background(), "test-voice-id", req, nil)
	if err != nil {
		t.Fatalf("ConvertWithTimestamps failed: %v", err)
	}

	// Verify response
	if result.AudioBase64 != "bW9jayBhdWRpbyBkYXRh" {
		t.Errorf("unexpected audio data")
	}
	if len(result.Alignment.Characters) != 5 {
		t.Errorf("expected 5 characters, got %d", len(result.Alignment.Characters))
	}
}

func TestConvertStream(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if !strings.Contains(r.URL.Path, "/stream") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Return mock streaming audio data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("streaming audio chunk 1"))
		w.Write([]byte("streaming audio chunk 2"))
	}))
	defer server.Close()

	// Create client with test server
	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	// Test request
	req := ConvertRequest{
		Text: "Hello world",
	}

	result, err := service.ConvertStream(context.Background(), "test-voice-id", req, nil)
	if err != nil {
		t.Fatalf("ConvertStream failed: %v", err)
	}
	defer result.Close()

	// Read response
	data, err := io.ReadAll(result)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	expected := "streaming audio chunk 1streaming audio chunk 2"
	if string(data) != expected {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestConvertError(t *testing.T) {
	// Create test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"detail": "Invalid voice ID"}`))
	}))
	defer server.Close()

	// Create client with test server
	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	// Test request
	req := ConvertRequest{
		Text: "Hello",
	}

	_, err := service.Convert(context.Background(), "invalid-voice-id", req, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verify error contains expected message
	if !strings.Contains(err.Error(), "Invalid voice ID") && !strings.Contains(err.Error(), "400") {
		t.Errorf("unexpected error message: %v", err)
	}
}
