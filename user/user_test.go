package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

func TestGetSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/user/subscription" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"tier": "free",
			"character_count": 5000,
			"character_limit": 10000,
			"can_extend_character_limit": false,
			"allowed_to_extend_character_limit": false,
			"next_character_count_reset_unix": 1738356858,
			"voice_limit": 10,
			"max_voice_add_edits": 50,
			"voice_add_edit_counter": 5,
			"professional_voice_limit": 0,
			"can_extend_voice_limit": false,
			"can_use_instant_voice_cloning": true,
			"can_use_professional_voice_cloning": false,
			"currency": "usd",
			"status": "free"
		}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.GetSubscription(context.Background())
	if err != nil {
		t.Fatalf("GetSubscription failed: %v", err)
	}

	if result.Tier != "free" {
		t.Errorf("unexpected tier: %s", result.Tier)
	}
	if result.CharacterCount != 5000 {
		t.Errorf("unexpected character count: %d", result.CharacterCount)
	}
	if !result.CanUseInstantVoiceCloning {
		t.Error("expected instant voice cloning to be available")
	}
}

func TestGetInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/user" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"subscription": {
				"tier": "free",
				"character_count": 5000,
				"character_limit": 10000,
				"can_extend_character_limit": false,
				"allowed_to_extend_character_limit": false,
				"next_character_count_reset_unix": 1738356858,
				"voice_limit": 10,
				"max_voice_add_edits": 50,
				"voice_add_edit_counter": 5,
				"professional_voice_limit": 0,
				"can_extend_voice_limit": false,
				"can_use_instant_voice_cloning": true,
				"can_use_professional_voice_cloning": false,
				"currency": "usd",
				"status": "free"
			},
			"is_new_user": false,
			"xi_api_key": "test-api-key",
			"can_use_delayed_payment_methods": false,
			"is_onboarding_completed": true,
			"is_onboarding_checklist_completed": true
		}`))
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-api-key", elevenlabs.WithBaseURL(server.URL))
	service := NewService(client)

	result, err := service.GetInfo(context.Background())
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

	if result.Subscription.Tier != "free" {
		t.Errorf("unexpected tier: %s", result.Subscription.Tier)
	}
	if result.XIAPIKey != "test-api-key" {
		t.Errorf("unexpected API key: %s", result.XIAPIKey)
	}
	if result.IsNewUser {
		t.Error("expected user to not be new")
	}
}
