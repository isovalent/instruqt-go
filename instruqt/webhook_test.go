package instruqt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/svix/svix-webhooks/go"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHandleWebhook tests the HandleWebhook function for different scenarios
func TestHandleWebhook(t *testing.T) {
	secret := "dGVzdC1zZWNyZXQ=" // Raw secret for testing

	handler := func(w http.ResponseWriter, r *http.Request, webhook WebhookEvent) error {
		// Custom processing for the test
		if webhook.Type == "test_event" {
			w.WriteHeader(http.StatusOK)
			return nil
		}
		return http.ErrNotSupported
	}

	webhookHandler := HandleWebhook(handler, secret)

	t.Run("Valid Webhook", func(t *testing.T) {
		webhookEvent := WebhookEvent{
			Type:          "test_event",
			TrackId:       "track_123",
			ParticipantId: "participant_123",
			UserId:        "user_123",
			Timestamp:     time.Now(),
		}
		payload, _ := json.Marshal(webhookEvent)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		// Add signature header
		wh, _ := svix.NewWebhook(secret)
		ts := time.Now()
		signature, _ := wh.Sign("test", ts, payload)
		req.Header.Set("Svix-Id", "test")
		req.Header.Set("Svix-Signature", signature)
		req.Header.Set("Svix-Timestamp", fmt.Sprintf("%v", ts.Unix()))
		req.Header.Set("Webhook-Id", "test")
		req.Header.Set("Webhook-Signature", signature)
		req.Header.Set("Webhook-Timestamp", fmt.Sprintf("%v", ts.Unix()))

		rr := httptest.NewRecorder()
		webhookHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Invalid Webhook Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
		rr := httptest.NewRecorder()
		webhookHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("Invalid Webhook Signature", func(t *testing.T) {
		ts := time.Now()
		webhookEvent := WebhookEvent{
			Type:          "test_event",
			TrackId:       "track_123",
			ParticipantId: "participant_123",
			UserId:        "user_123",
			Timestamp:     ts,
		}
		payload, _ := json.Marshal(webhookEvent)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		// Add an invalid signature header
		req.Header.Set("Svix-Id", "test")
		req.Header.Set("Svix-Signature", "invalid-signature")
		req.Header.Set("Svix-Timestamp", fmt.Sprintf("%v", ts.Unix()))
		req.Header.Set("Webhook-Id", "test")
		req.Header.Set("Webhook-Signature", "invalid-signature")
		req.Header.Set("Webhook-Timestamp", fmt.Sprintf("%v", ts.Unix()))

		rr := httptest.NewRecorder()
		webhookHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}
