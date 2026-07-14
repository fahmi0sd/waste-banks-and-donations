package pkg

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestMailjetClient(baseURL string) *MailjetClient {
	return NewMailjetClient("test-api-key", "test-secret-key", "sender@banksampah.test", "Bank Sampah", baseURL)
}

func TestMailjetClient_SendEmail_Success(t *testing.T) {
	var gotAuthHeader string
	var gotBody mailjetRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestMailjetClient(server.URL)

	err := client.SendEmail("user@example.com", "Budi", "Deposit Bertambah", "Halo Budi, deposit kamu bertambah.")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !strings.HasPrefix(gotAuthHeader, "Basic ") {
		t.Errorf("expected Basic auth header, got: %q", gotAuthHeader)
	}

	if len(gotBody.Messages) != 1 {
		t.Fatalf("expected 1 message, got: %d", len(gotBody.Messages))
	}
	msg := gotBody.Messages[0]
	if msg.From.Email != "sender@banksampah.test" {
		t.Errorf("expected sender sender@banksampah.test, got: %s", msg.From.Email)
	}
	if len(msg.To) != 1 || msg.To[0].Email != "user@example.com" {
		t.Errorf("expected recipient user@example.com, got: %+v", msg.To)
	}
	if msg.Subject != "Deposit Bertambah" {
		t.Errorf("expected subject 'Deposit Bertambah', got: %q", msg.Subject)
	}
}

func TestMailjetClient_SendEmail_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized) // simulasi API key salah
	}))
	defer server.Close()

	client := newTestMailjetClient(server.URL)

	err := client.SendEmail("user@example.com", "Budi", "Subject", "Content")
	if err == nil {
		t.Fatal("expected error when Mailjet returns non-2xx status, got nil")
	}
}

func TestMailjetClient_SendEmail_NetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	client := newTestMailjetClient(server.URL)

	err := client.SendEmail("user@example.com", "Budi", "Subject", "Content")
	if err == nil {
		t.Fatal("expected error when server is unreachable, got nil")
	}
}
