package voices

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

func TestList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/voices" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"voices": [
				{
					"voice_id": "21m00Tcm4TlvDq8ikWAM",
					"name": "Rachel",
					"category": "premade"
				}
			]
		}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(result.Voices) != 1 {
		t.Errorf("expected 1 voice, got %d", len(result.Voices))
	}
	if result.Voices[0].Name != "Rachel" {
		t.Errorf("unexpected voice name: %s", result.Voices[0].Name)
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v1/voices/test-voice-id") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"voice_id": "test-voice-id",
			"name": "Test Voice",
			"category": "cloned"
		}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.Get(context.Background(), "test-voice-id")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if result.VoiceID != "test-voice-id" {
		t.Errorf("unexpected voice ID: %s", result.VoiceID)
	}
	if result.Name != "Test Voice" {
		t.Errorf("unexpected voice name: %s", result.Name)
	}
}

func TestGetDefaultSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/voices/settings/default" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"stability": 0.5,
			"similarity_boost": 0.75,
			"style": 0,
			"use_speaker_boost": true
		}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.GetDefaultSettings(context.Background())
	if err != nil {
		t.Fatalf("GetDefaultSettings failed: %v", err)
	}

	if result.Stability != 0.5 {
		t.Errorf("unexpected stability: %f", result.Stability)
	}
	if result.SimilarityBoost != 0.75 {
		t.Errorf("unexpected similarity boost: %f", result.SimilarityBoost)
	}
}

func TestGetSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/settings") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"stability": 0.6,
			"similarity_boost": 0.8
		}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.GetSettings(context.Background(), "test-voice-id")
	if err != nil {
		t.Fatalf("GetSettings failed: %v", err)
	}

	if result.Stability != 0.6 {
		t.Errorf("unexpected stability: %f", result.Stability)
	}
}

func TestUpdateSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/settings/edit") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	settings := VoiceSettings{
		Stability:       0.7,
		SimilarityBoost: 0.85,
	}

	err := service.UpdateSettings(context.Background(), "test-voice-id", settings)
	if err != nil {
		t.Fatalf("UpdateSettings failed: %v", err)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/v1/voices/test-voice-id") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	err := service.Delete(context.Background(), "test-voice-id")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/edit") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	req := UpdateVoiceRequest{
		Name:        "Updated Voice",
		Description: "Updated description",
	}

	err := service.Update(context.Background(), "test-voice-id", req)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
}
