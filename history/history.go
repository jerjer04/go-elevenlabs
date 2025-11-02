package history

import (
	"context"
	"fmt"
	"io"
	"strconv"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

// Service handles history operations
type Service struct {
	client *elevenlabs.Client
}

// NewService creates a new history service
func NewService(client *elevenlabs.Client) *Service {
	return &Service{client: client}
}

// HistoryItem represents a generated audio history item
type HistoryItem struct {
	HistoryItemID            string                 `json:"history_item_id"`
	RequestID                string                 `json:"request_id"`
	VoiceID                  string                 `json:"voice_id"`
	ModelID                  string                 `json:"model_id"`
	VoiceName                string                 `json:"voice_name"`
	VoiceCategory            string                 `json:"voice_category"`
	Text                     string                 `json:"text"`
	DateUnix                 int64                  `json:"date_unix"`
	CharacterCountChangeFrom int                    `json:"character_count_change_from"`
	CharacterCountChangeTo   int                    `json:"character_count_change_to"`
	ContentType              string                 `json:"content_type"`
	State                    string                 `json:"state"`
	Settings                 map[string]interface{} `json:"settings,omitempty"`
	Feedback                 *Feedback              `json:"feedback,omitempty"`
	ShareLinkID              string                 `json:"share_link_id,omitempty"`
	Source                   string                 `json:"source,omitempty"`
	Alignments               *Alignments            `json:"alignments,omitempty"`
}

// Feedback represents user feedback on a history item
type Feedback struct {
	ThumbsUp        bool   `json:"thumbs_up"`
	Feedback        string `json:"feedback,omitempty"`
	Emotions        bool   `json:"emotions,omitempty"`
	InaccurateClone bool   `json:"inaccurate_clone,omitempty"`
	Glitches        bool   `json:"glitches,omitempty"`
	AudioQuality    bool   `json:"audio_quality,omitempty"`
	Other           bool   `json:"other,omitempty"`
	ReviewStatus    string `json:"review_status,omitempty"`
}

// Alignments contains character timing alignments
type Alignments struct {
	Alignment           *Alignment `json:"alignment,omitempty"`
	NormalizedAlignment *Alignment `json:"normalized_alignment,omitempty"`
}

// Alignment contains character-level timing information
type Alignment struct {
	Characters              []string  `json:"characters"`
	CharacterStartTimesSecs []float64 `json:"character_start_times_seconds"`
	CharacterEndTimesSecs   []float64 `json:"character_end_times_seconds"`
}

// ListHistoryResponse contains a list of history items
type ListHistoryResponse struct {
	History           []HistoryItem `json:"history"`
	LastHistoryItemID string        `json:"last_history_item_id"`
	HasMore           bool          `json:"has_more"`
}

// ListOptions contains optional parameters for listing history
type ListOptions struct {
	PageSize int    // Maximum number of items to return (default: 100, max: 1000)
	VoiceID  string // Filter by voice ID
}

// List retrieves generated audio history items
func (s *Service) List(ctx context.Context, opts *ListOptions) (*ListHistoryResponse, error) {
	path := "/v1/history"

	if opts != nil {
		params := make(map[string]string)
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
		if opts.VoiceID != "" {
			params["voice_id"] = opts.VoiceID
		}
		path = elevenlabs.BuildURL(path, params)
	}

	var result ListHistoryResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a specific history item by ID
func (s *Service) Get(ctx context.Context, historyItemID string) (*HistoryItem, error) {
	path := fmt.Sprintf("/v1/history/%s", historyItemID)
	var result HistoryItem
	if err := s.client.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a history item by ID
func (s *Service) Delete(ctx context.Context, historyItemID string) error {
	path := fmt.Sprintf("/v1/history/%s", historyItemID)
	var result map[string]interface{}
	return s.client.DoRequest(ctx, "DELETE", path, nil, &result)
}

// GetAudio retrieves the audio for a specific history item
// Returns an io.ReadCloser that must be closed by the caller
func (s *Service) GetAudio(ctx context.Context, historyItemID string) (io.ReadCloser, error) {
	path := fmt.Sprintf("/v1/history/%s/audio", historyItemID)
	return s.client.DoRequestStream(ctx, "GET", path, nil)
}

// DownloadRequest contains parameters for downloading history items
type DownloadRequest struct {
	HistoryItemIDs []string `json:"history_item_ids"`
}

// Download downloads one or more history items
// If one ID is provided, returns a single audio file
// If multiple IDs are provided, returns a ZIP file
// Returns an io.ReadCloser that must be closed by the caller
func (s *Service) Download(ctx context.Context, historyItemIDs []string) (io.ReadCloser, error) {
	req := DownloadRequest{
		HistoryItemIDs: historyItemIDs,
	}
	return s.client.DoRequestStream(ctx, "POST", "/v1/history/download", req)
}
