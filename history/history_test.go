package history

import (
	"context"
	"io"
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
		if r.URL.Path != "/v1/history" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"history": [
				{
					"history_item_id": "test-id",
					"voice_name": "Rachel",
					"text": "Hello world",
					"date_unix": 1234567890
				}
			],
			"last_history_item_id": "test-id",
			"has_more": false
		}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(result.History) != 1 {
		t.Errorf("expected 1 history item, got %d", len(result.History))
	}
	if result.History[0].VoiceName != "Rachel" {
		t.Errorf("unexpected voice name: %s", result.History[0].VoiceName)
	}
}

func TestListWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageSize := r.URL.Query().Get("page_size")
		if pageSize != "10" {
			t.Errorf("expected page_size=10, got %s", pageSize)
		}

		voiceID := r.URL.Query().Get("voice_id")
		if voiceID != "test-voice" {
			t.Errorf("expected voice_id=test-voice, got %s", voiceID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"history": [], "last_history_item_id": "", "has_more": false}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	opts := &ListOptions{
		PageSize: 10,
		VoiceID:  "test-voice",
	}

	_, err := service.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v1/history/test-item-id") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"history_item_id": "test-item-id",
			"voice_name": "Rachel",
			"text": "Test text",
			"date_unix": 1234567890
		}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.Get(context.Background(), "test-item-id")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if result.HistoryItemID != "test-item-id" {
		t.Errorf("unexpected history item ID: %s", result.HistoryItemID)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/v1/history/test-item-id") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	err := service.Delete(context.Background(), "test-item-id")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestGetAudio(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/audio") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock audio data"))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.GetAudio(context.Background(), "test-item-id")
	if err != nil {
		t.Fatalf("GetAudio failed: %v", err)
	}
	defer result.Close()

	data, err := io.ReadAll(result)
	if err != nil {
		t.Fatalf("failed to read audio: %v", err)
	}

	if string(data) != "mock audio data" {
		t.Errorf("unexpected audio data: %s", string(data))
	}
}

func TestDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/download") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock download data"))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.Download(context.Background(), []string{"item1", "item2"})
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	defer result.Close()

	data, err := io.ReadAll(result)
	if err != nil {
		t.Fatalf("failed to read download: %v", err)
	}

	if string(data) != "mock download data" {
		t.Errorf("unexpected download data: %s", string(data))
	}
}
