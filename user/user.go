package user

import (
	"context"

	elevenlabs "github.com/dhia-gharsallaoui/go-elevenlabs"
)

// Service handles user operations
type Service struct {
	client *elevenlabs.Client
}

// NewService creates a new user service
func NewService(client *elevenlabs.Client) *Service {
	return &Service{client: client}
}

// Subscription contains user subscription information
type Subscription struct {
	Tier                           string       `json:"tier"`
	CharacterCount                 int          `json:"character_count"`
	CharacterLimit                 int          `json:"character_limit"`
	CanExtendCharacterLimit        bool         `json:"can_extend_character_limit"`
	AllowedToExtendCharacterLimit  bool         `json:"allowed_to_extend_character_limit"`
	NextCharacterCountResetUnix    int64        `json:"next_character_count_reset_unix"`
	VoiceLimit                     int          `json:"voice_limit"`
	MaxVoiceAddEdits               int          `json:"max_voice_add_edits"`
	VoiceAddEditCounter            int          `json:"voice_add_edit_counter"`
	ProfessionalVoiceLimit         int          `json:"professional_voice_limit"`
	CanExtendVoiceLimit            bool         `json:"can_extend_voice_limit"`
	CanUseInstantVoiceCloning      bool         `json:"can_use_instant_voice_cloning"`
	CanUseProfessionalVoiceCloning bool         `json:"can_use_professional_voice_cloning"`
	Currency                       string       `json:"currency"`
	Status                         string       `json:"status"`
	BillingPeriod                  string       `json:"billing_period,omitempty"`
	CharacterRefreshPeriod         string       `json:"character_refresh_period,omitempty"`
	NextInvoice                    *NextInvoice `json:"next_invoice,omitempty"`
	HasOpenInvoices                bool         `json:"has_open_invoices,omitempty"`
}

// NextInvoice contains information about the next invoice
type NextInvoice struct {
	AmountDueCents         int   `json:"amount_due_cents"`
	NextPaymentAttemptUnix int64 `json:"next_payment_attempt_unix"`
}

// User contains user account information
type User struct {
	Subscription                   *Subscription `json:"subscription"`
	IsNewUser                      bool          `json:"is_new_user"`
	XIAPIKey                       string        `json:"xi_api_key"`
	CanUseDelayedPaymentMethods    bool          `json:"can_use_delayed_payment_methods"`
	IsOnboardingCompleted          bool          `json:"is_onboarding_completed"`
	IsOnboardingChecklistCompleted bool          `json:"is_onboarding_checklist_completed"`
	FirstName                      string        `json:"first_name,omitempty"`
	IsAPIKeyHashed                 bool          `json:"is_api_key_hashed,omitempty"`
	XIAPIKeyPreview                string        `json:"xi_api_key_preview,omitempty"`
	ReferralLinkCode               string        `json:"referral_link_code,omitempty"`
	PartnerStackPartnerDefaultLink string        `json:"partnerstack_partner_default_link,omitempty"`
}

// GetSubscription retrieves user subscription information
func (s *Service) GetSubscription(ctx context.Context) (*Subscription, error) {
	var result Subscription
	if err := s.client.DoRequest(ctx, "GET", "/v1/user/subscription", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInfo retrieves user account information
func (s *Service) GetInfo(ctx context.Context) (*User, error) {
	var result User
	if err := s.client.DoRequest(ctx, "GET", "/v1/user", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
