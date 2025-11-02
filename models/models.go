package models

import (
	"context"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

// Service handles model operations
type Service struct {
	client *elevenlabs.Client
}

// NewService creates a new models service
func NewService(client *elevenlabs.Client) *Service {
	return &Service{client: client}
}

// Model represents an ElevenLabs TTS model
type Model struct {
	ModelID                   string   `json:"model_id"`
	Name                      string   `json:"name"`
	CanBeFinetuned            bool     `json:"can_be_finetuned"`
	CanDoTextToSpeech         bool     `json:"can_do_text_to_speech"`
	CanDoVoiceConversion      bool     `json:"can_do_voice_conversion"`
	CanUseStyle               bool     `json:"can_use_style"`
	CanUseSpeakerBoost        bool     `json:"can_use_speaker_boost"`
	ServesProVoices           bool     `json:"serves_pro_voices"`
	TokenCostFactor           float64  `json:"token_cost_factor"`
	Description               string   `json:"description,omitempty"`
	RequiresAlphaAccess       bool     `json:"requires_alpha_access,omitempty"`
	MaxCharactersRequestFree  int      `json:"max_characters_request_free,omitempty"`
	MaxCharactersRequestPaid  int      `json:"max_characters_request_paid,omitempty"`
	Languages                 []Language `json:"languages,omitempty"`
}

// Language represents a supported language
type Language struct {
	LanguageID string `json:"language_id"`
	Name       string `json:"name"`
}

// ListModelsResponse contains a list of available models
type ListModelsResponse struct {
	Models []Model `json:"models,omitempty"`
}

// List retrieves all available TTS models
func (s *Service) List(ctx context.Context) (*ListModelsResponse, error) {
	var result []Model
	if err := s.client.DoRequest(ctx, "GET", "/v1/models", nil, &result); err != nil {
		return nil, err
	}
	return &ListModelsResponse{Models: result}, nil
}
