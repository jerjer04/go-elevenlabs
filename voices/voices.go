package voices

import (
	"context"
	"fmt"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

// Service handles voice operations
type Service struct {
	client *elevenlabs.Client
}

// NewService creates a new voices service
func NewService(client *elevenlabs.Client) *Service {
	return &Service{client: client}
}

// VoiceSettings contains voice configuration parameters
type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style,omitempty"`
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

// Voice represents a voice in the ElevenLabs system
type Voice struct {
	VoiceID                  string                `json:"voice_id"`
	Name                     string                `json:"name"`
	Samples                  []Sample              `json:"samples,omitempty"`
	Category                 string                `json:"category"`
	FineTuning               *FineTuning           `json:"fine_tuning,omitempty"`
	Labels                   map[string]string     `json:"labels,omitempty"`
	Description              string                `json:"description,omitempty"`
	PreviewURL               string                `json:"preview_url,omitempty"`
	AvailableForTiers        []string              `json:"available_for_tiers,omitempty"`
	Settings                 *VoiceSettings        `json:"settings,omitempty"`
	Sharing                  *SharingSettings      `json:"sharing,omitempty"`
	HighQualityBaseModelIDs  []string              `json:"high_quality_base_model_ids,omitempty"`
	SafetyControl            string                `json:"safety_control,omitempty"`
	VoiceVerification        *VoiceVerification    `json:"voice_verification,omitempty"`
	PermissionOnResource     string                `json:"permission_on_resource,omitempty"`
	IsOwner                  bool                  `json:"is_owner,omitempty"`
	IsLegacy                 bool                  `json:"is_legacy,omitempty"`
	IsMixed                  bool                  `json:"is_mixed,omitempty"`
	CreatedAtUnix            int64                 `json:"created_at_unix,omitempty"`
}

// Sample represents an audio sample for a voice
type Sample struct {
	SampleID  string `json:"sample_id"`
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	SizeBytes int    `json:"size_bytes"`
	Hash      string `json:"hash"`
}

// FineTuning contains voice fine-tuning information
type FineTuning struct {
	IsAllowedToFineTune               bool                    `json:"is_allowed_to_fine_tune"`
	FinetuningState                   string                  `json:"finetuning_state,omitempty"`
	VerificationFailures              []string                `json:"verification_failures,omitempty"`
	VerificationAttemptsCount         int                     `json:"verification_attempts_count"`
	ManualVerificationRequested       bool                    `json:"manual_verification_requested"`
	Language                          string                  `json:"language,omitempty"`
	Progress                          map[string]interface{}  `json:"progress,omitempty"`
	Message                           map[string]interface{}  `json:"message,omitempty"`
	DatasetDurationSeconds            float64                 `json:"dataset_duration_seconds,omitempty"`
	VerificationAttempts              []VerificationAttempt   `json:"verification_attempts,omitempty"`
	SliceIDs                          []string                `json:"slice_ids,omitempty"`
	ManualVerification                *ManualVerification     `json:"manual_verification,omitempty"`
	MaxVerificationAttempts           int                     `json:"max_verification_attempts,omitempty"`
	NextMaxVerificationAttemptsResetUnixMs int64              `json:"next_max_verification_attempts_reset_unix_ms,omitempty"`
}

// VerificationAttempt represents a voice verification attempt
type VerificationAttempt struct {
	Text                 string     `json:"text"`
	DateUnix             int64      `json:"date_unix"`
	Accepted             bool       `json:"accepted"`
	Similarity           float64    `json:"similarity"`
	LevenshteinDistance  float64    `json:"levenshtein_distance"`
	Recording            *Recording `json:"recording,omitempty"`
}

// Recording represents an audio recording
type Recording struct {
	RecordingID      string `json:"recording_id"`
	MimeType         string `json:"mime_type"`
	SizeBytes        int    `json:"size_bytes"`
	UploadDateUnix   int64  `json:"upload_date_unix"`
	Transcription    string `json:"transcription,omitempty"`
}

// ManualVerification contains manual verification information
type ManualVerification struct {
	ExtraText       string `json:"extra_text,omitempty"`
	RequestTimeUnix int64  `json:"request_time_unix,omitempty"`
	Files           []File `json:"files,omitempty"`
}

// File represents an uploaded file
type File struct {
	FileID         string `json:"file_id"`
	FileName       string `json:"file_name"`
	MimeType       string `json:"mime_type"`
	SizeBytes      int    `json:"size_bytes"`
	UploadDateUnix int64  `json:"upload_date_unix"`
}

// SharingSettings contains voice sharing configuration
type SharingSettings struct {
	Status                   string            `json:"status"`
	HistoryItemSampleID      string            `json:"history_item_sample_id,omitempty"`
	DateUnix                 int64             `json:"date_unix,omitempty"`
	WhitelistedEmails        []string          `json:"whitelisted_emails,omitempty"`
	PublicOwnerID            string            `json:"public_owner_id,omitempty"`
	OriginalVoiceID          string            `json:"original_voice_id,omitempty"`
	FinancialRewardsEnabled  bool              `json:"financial_rewards_enabled,omitempty"`
	FreeUsersAllowed         bool              `json:"free_users_allowed,omitempty"`
	LiveModerationEnabled    bool              `json:"live_moderation_enabled,omitempty"`
	Rate                     float64           `json:"rate,omitempty"`
	NoticePeriod             int               `json:"notice_period,omitempty"`
	DisableAtUnix            int64             `json:"disable_at_unix,omitempty"`
	VoiceMixingAllowed       bool              `json:"voice_mixing_allowed,omitempty"`
	Featured                 bool              `json:"featured,omitempty"`
	Category                 string            `json:"category,omitempty"`
	ReaderAppEnabled         bool              `json:"reader_app_enabled,omitempty"`
	ImageURL                 string            `json:"image_url,omitempty"`
	BanReason                string            `json:"ban_reason,omitempty"`
	LikedByCount             int               `json:"liked_by_count,omitempty"`
	ClonedByCount            int               `json:"cloned_by_count,omitempty"`
	Name                     string            `json:"name,omitempty"`
	Description              string            `json:"description,omitempty"`
	Labels                   map[string]string `json:"labels,omitempty"`
	ReviewStatus             string            `json:"review_status,omitempty"`
	ReviewMessage            string            `json:"review_message,omitempty"`
	EnabledInLibrary         bool              `json:"enabled_in_library,omitempty"`
	InstagramUsername        string            `json:"instagram_username,omitempty"`
	TwitterUsername          string            `json:"twitter_username,omitempty"`
	YouTubeUsername          string            `json:"youtube_username,omitempty"`
	TikTokUsername           string            `json:"tiktok_username,omitempty"`
	ModerationCheck          *ModerationCheck  `json:"moderation_check,omitempty"`
	ReaderRestrictedOn       []ResourceRestriction `json:"reader_restricted_on,omitempty"`
}

// ModerationCheck contains moderation check results
type ModerationCheck struct {
	DateCheckedUnix     int64     `json:"date_checked_unix"`
	NameValue           string    `json:"name_value"`
	NameCheck           bool      `json:"name_check"`
	DescriptionValue    string    `json:"description_value"`
	DescriptionCheck    bool      `json:"description_check"`
	SampleIDs           []string  `json:"sample_ids,omitempty"`
	SampleChecks        []float64 `json:"sample_checks,omitempty"`
	CaptchaIDs          []string  `json:"captcha_ids,omitempty"`
	CaptchaChecks       []float64 `json:"captcha_checks,omitempty"`
}

// ResourceRestriction represents a resource restriction
type ResourceRestriction struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
}

// VoiceVerification contains voice verification information
type VoiceVerification struct {
	RequiresVerification      bool                  `json:"requires_verification"`
	IsVerified                bool                  `json:"is_verified"`
	VerificationFailures      []string              `json:"verification_failures,omitempty"`
	VerificationAttemptsCount int                   `json:"verification_attempts_count"`
	Language                  string                `json:"language,omitempty"`
	VerificationAttempts      []VerificationAttempt `json:"verification_attempts,omitempty"`
}

// ListVoicesResponse contains a list of voices
type ListVoicesResponse struct {
	Voices []Voice `json:"voices"`
}

// List retrieves all available voices
func (s *Service) List(ctx context.Context) (*ListVoicesResponse, error) {
	var result ListVoicesResponse
	if err := s.client.DoRequest(ctx, "GET", "/v1/voices", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves details for a specific voice
func (s *Service) Get(ctx context.Context, voiceID string) (*Voice, error) {
	path := fmt.Sprintf("/v1/voices/%s", voiceID)
	var result Voice
	if err := s.client.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDefaultSettings retrieves the default voice settings
func (s *Service) GetDefaultSettings(ctx context.Context) (*VoiceSettings, error) {
	var result VoiceSettings
	if err := s.client.DoRequest(ctx, "GET", "/v1/voices/settings/default", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSettings retrieves settings for a specific voice
func (s *Service) GetSettings(ctx context.Context, voiceID string) (*VoiceSettings, error) {
	path := fmt.Sprintf("/v1/voices/%s/settings", voiceID)
	var result VoiceSettings
	if err := s.client.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateSettings updates settings for a specific voice
func (s *Service) UpdateSettings(ctx context.Context, voiceID string, settings VoiceSettings) error {
	path := fmt.Sprintf("/v1/voices/%s/settings/edit", voiceID)
	return s.client.DoRequest(ctx, "POST", path, settings, nil)
}

// Delete removes a voice
func (s *Service) Delete(ctx context.Context, voiceID string) error {
	path := fmt.Sprintf("/v1/voices/%s", voiceID)
	var result map[string]interface{}
	return s.client.DoRequest(ctx, "DELETE", path, nil, &result)
}

// AddVoiceRequest contains parameters for adding a new voice
type AddVoiceRequest struct {
	Name        string            `json:"name"`
	Files       []byte            `json:"files"` // Multipart form data
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// UpdateVoiceRequest contains parameters for updating a voice
type UpdateVoiceRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Update updates a voice's metadata
func (s *Service) Update(ctx context.Context, voiceID string, req UpdateVoiceRequest) error {
	path := fmt.Sprintf("/v1/voices/%s/edit", voiceID)
	return s.client.DoRequest(ctx, "POST", path, req, nil)
}

// SharedVoicesResponse contains a list of shared voices
type SharedVoicesResponse struct {
	Voices []Voice `json:"voices"`
}

// GetSharedVoices retrieves voices shared with the user
func (s *Service) GetSharedVoices(ctx context.Context) (*SharedVoicesResponse, error) {
	var result SharedVoicesResponse
	if err := s.client.DoRequest(ctx, "GET", "/v1/shared-voices", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SimilarVoicesRequest contains parameters for finding similar voices
type SimilarVoicesRequest struct {
	AudioFile []byte `json:"audio_file"` // Multipart form data
}

// GetSimilarLibraryVoices finds voices similar to the provided audio
func (s *Service) GetSimilarLibraryVoices(ctx context.Context, req SimilarVoicesRequest) (*ListVoicesResponse, error) {
	var result ListVoicesResponse
	if err := s.client.DoRequest(ctx, "POST", "/v1/similar-voices", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
