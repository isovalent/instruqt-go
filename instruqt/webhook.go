package instruqt

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	svix "github.com/svix/svix-webhooks/go"
)

// WebhookEvent represents the structure of an incoming webhook event from Instruqt
// This includes details about the type of event, participant, and other metadata
// related to challenges, reviews, and custom parameters.
type WebhookEvent struct {
	Type             string            `json:"type"`              // Type of the event (e.g., challenge.completed)
	TrackId          string            `json:"track_id"`          // ID of the track related to the event
	TrackSlug        string            `json:"track_slug"`        // Slug identifier for the track
	ParticipantId    string            `json:"participant_id"`    // ID of the participant
	UserId           string            `json:"user_id"`           // ID of the user who triggered the event
	InviteId         string            `json:"invite_id"`         // ID of the invite associated with the event
	ClaimId          string            `json:"claim_id"`          // Claim ID associated with the event
	Timestamp        time.Time         `json:"timestamp"`         // Timestamp when the event occurred
	Reason           string            `json:"reason"`            // Reason for the event, if applicable
	Duration         int               `json:"duration"`          // Duration of the activity, if applicable
	CustomParameters map[string]string `json:"custom_parameters"` // Custom parameters for the sandbox

	// Challenges
	ChallengeId     string `json:"challenge_id"`     // ID of the challenge
	ChallengeIndex  int    `json:"challenge_index"`  // Index of the challenge in the track
	TotalChallenges int    `json:"total_challenges"` // Total number of challenges in the track

	// Review
	Content  string `json:"content"`   // Content of the review, if the event is related to a review
	ReviewId string `json:"review_id"` // ID of the review
	Score    int    `json:"score"`     // Score given in the review
}

// WebhookHandler is a handler function for processing webhooks
// It takes an HTTP response writer, an HTTP request, and a WebhookEvent structure
// and returns an error if the processing fails.
type WebhookHandler func(w http.ResponseWriter, r *http.Request, webhook WebhookEvent) error

// HandleWebhook is an HTTP handler that validates and processes incoming webhooks
// It takes a WebhookHandler function and a secret for validating the webhook signature.
func HandleWebhook(handler WebhookHandler, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		wh, err := svix.NewWebhook(secret)
		if err != nil {
			http.Error(w, "Failed to create webhook validator", http.StatusInternalServerError)
			return
		}

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "No payload", http.StatusBadRequest)
			return
		}

		err = wh.Verify(payload, r.Header)
		if err != nil {
			http.Error(w, "Invalid webhook signature", http.StatusUnauthorized)
			return
		}

		var webhook WebhookEvent
		if err := json.Unmarshal(payload, &webhook); err != nil {
			http.Error(w, "Failed to decode webhook payload", http.StatusBadRequest)
			return
		}

		if webhook.Type == "" {
			http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
			return
		}

		if err := handler(w, r, webhook); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
