package tts

import (
	"context"
	"fmt"
	"io"
	"strconv"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

// Service handles text-to-speech operations
type Service struct {
	client *elevenlabs.Client
}

// NewService creates a new TTS service
func NewService(client *elevenlabs.Client) *Service {
	return &Service{client: client}
}

// VoiceSettings contains voice configuration parameters
type VoiceSettings struct {
	Stability       float64 `json:"stability,omitempty"`
	SimilarityBoost float64 `json:"similarity_boost,omitempty"`
	Style           float64 `json:"style,omitempty"`
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

// ConvertRequest contains parameters for text-to-speech conversion
type ConvertRequest struct {
	Text                            string                     `json:"text"`
	ModelID                         string                     `json:"model_id,omitempty"`
	VoiceSettings                   *VoiceSettings             `json:"voice_settings,omitempty"`
	PronunciationDictionaryLocators []PronunciationDictLocator `json:"pronunciation_dictionary_locators,omitempty"`
}

// PronunciationDictLocator specifies a pronunciation dictionary to use
type PronunciationDictLocator struct {
	PronunciationDictionaryID string `json:"pronunciation_dictionary_id"`
	VersionID                 string `json:"version_id"`
}

// ConvertOptions contains optional parameters for TTS conversion
type ConvertOptions struct {
	OptimizeStreamingLatency int    // Latency optimization level (0-4)
	OutputFormat             string // Audio output format (e.g., "mp3_22050_32", "mp3_44100_128")
}

// WithTimestampsResponse contains audio data with character timing information
type WithTimestampsResponse struct {
	AudioBase64         string     `json:"audio_base64"`
	Alignment           Alignment  `json:"alignment"`
	NormalizedAlignment *Alignment `json:"normalized_alignment,omitempty"`
}

// Alignment contains character-level timing information
type Alignment struct {
	Characters              []string  `json:"characters"`
	CharacterStartTimesSecs []float64 `json:"character_start_times_seconds"`
	CharacterEndTimesSecs   []float64 `json:"character_end_times_seconds"`
}

// Convert performs text-to-speech conversion and returns audio data
// Returns an io.ReadCloser that must be closed by the caller
func (s *Service) Convert(ctx context.Context, voiceID string, req ConvertRequest, opts *ConvertOptions) (io.ReadCloser, error) {
	path := fmt.Sprintf("/v1/text-to-speech/%s", voiceID)

	// Add query parameters if options are provided
	if opts != nil {
		params := make(map[string]string)
		if opts.OptimizeStreamingLatency > 0 {
			params["optimize_streaming_latency"] = strconv.Itoa(opts.OptimizeStreamingLatency)
		}
		if opts.OutputFormat != "" {
			params["output_format"] = opts.OutputFormat
		}
		path = elevenlabs.BuildURL(path, params)
	}

	return s.client.DoRequestStream(ctx, "POST", path, req)
}

// ConvertWithTimestamps performs text-to-speech conversion and returns audio with character-level timing
func (s *Service) ConvertWithTimestamps(ctx context.Context, voiceID string, req ConvertRequest, opts *ConvertOptions) (*WithTimestampsResponse, error) {
	path := fmt.Sprintf("/v1/text-to-speech/%s/with-timestamps", voiceID)

	// Add query parameters if options are provided
	if opts != nil {
		params := make(map[string]string)
		if opts.OptimizeStreamingLatency > 0 {
			params["optimize_streaming_latency"] = strconv.Itoa(opts.OptimizeStreamingLatency)
		}
		if opts.OutputFormat != "" {
			params["output_format"] = opts.OutputFormat
		}
		path = elevenlabs.BuildURL(path, params)
	}

	var result WithTimestampsResponse
	if err := s.client.DoRequest(ctx, "POST", path, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ConvertStream performs streaming text-to-speech conversion
// Returns an io.ReadCloser that must be closed by the caller
func (s *Service) ConvertStream(ctx context.Context, voiceID string, req ConvertRequest, opts *ConvertOptions) (io.ReadCloser, error) {
	path := fmt.Sprintf("/v1/text-to-speech/%s/stream", voiceID)

	// Add query parameters if options are provided
	if opts != nil {
		params := make(map[string]string)
		if opts.OptimizeStreamingLatency > 0 {
			params["optimize_streaming_latency"] = strconv.Itoa(opts.OptimizeStreamingLatency)
		}
		if opts.OutputFormat != "" {
			params["output_format"] = opts.OutputFormat
		}
		path = elevenlabs.BuildURL(path, params)
	}

	return s.client.DoRequestStream(ctx, "POST", path, req)
}

// ConvertStreamWithTimestamps performs streaming text-to-speech conversion with character-level timing
// Returns an io.ReadCloser that must be closed by the caller
func (s *Service) ConvertStreamWithTimestamps(ctx context.Context, voiceID string, req ConvertRequest, opts *ConvertOptions) (io.ReadCloser, error) {
	path := fmt.Sprintf("/v1/text-to-speech/%s/stream/with-timestamps", voiceID)

	// Add query parameters if options are provided
	if opts != nil {
		params := make(map[string]string)
		if opts.OptimizeStreamingLatency > 0 {
			params["optimize_streaming_latency"] = strconv.Itoa(opts.OptimizeStreamingLatency)
		}
		if opts.OutputFormat != "" {
			params["output_format"] = opts.OutputFormat
		}
		path = elevenlabs.BuildURL(path, params)
	}

	return s.client.DoRequestStream(ctx, "POST", path, req)
}
