//go:build integration
// +build integration

package gobark

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		key     string
		wantErr bool
	}{
		{
			name:    "valid client with default base URL",
			baseURL: "",
			key:     "test-key",
			wantErr: false,
		},
		{
			name:    "valid client with custom base URL",
			baseURL: "https://custom.bark.server",
			key:     "test-key",
			wantErr: false,
		},
		{
			name:    "invalid client without key",
			baseURL: "https://api.day.app",
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client without error")
			}
		})
	}
}

func TestBuildNotificationURL(t *testing.T) {
	client, _ := NewClient("https://api.day.app", "test-key")
	tests := []struct {
		name     string
		body     string
		opts     []Option
		wantPath string
	}{
		{
			name:     "simple notification",
			body:     "test message",
			opts:     nil,
			wantPath: "https://api.day.app/test-key/test%20message",
		},
		{
			name: "notification with title",
			body: "test message",
			opts: []Option{
				WithTitle("Test Title"),
			},
			wantPath: "https://api.day.app/test-key/Test%20Title/test%20message",
		},
		{
			name: "full notification",
			body: "test message",
			opts: []Option{
				WithTitle("Test Title"),
				WithSubtitle("Test Subtitle"),
				WithSound("bell"),
				WithTimeSensitive(),
			},
			wantPath: "https://api.day.app/test-key/Test%20Title/Test%20Subtitle/test%20message?level=timeSensitive&sound=bell",
		},
		{
			name: "notification with newlines",
			body: "Line 1\nLine 2\nLine 3",
			opts: []Option{
				WithTitle("Multiline Test"),
			},
			wantPath: "https://api.day.app/test-key/Multiline%20Test/Line%201%0ALine%202%0ALine%203",
		},
		{
			name: "notification with special characters",
			body: "Special chars: !@#$%^&*()",
			opts: []Option{
				WithTitle("Special Chars"),
			},
			wantPath: "https://api.day.app/test-key/Special%20Chars/Special%20chars%3A%20%21%40%23%24%25%5E%26%2A%28%29",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &notification{
				title: defaultTitle,
				body:  tt.body,
			}
			for _, opt := range tt.opts {
				opt(n)
			}

			urlPath := client.buildNotificationURL(n)

			// For query parameters, the order might be different, so we need to check differently
			if strings.Contains(tt.wantPath, "?") {
				parts := strings.Split(tt.wantPath, "?")
				basePath := parts[0]
				queryParams := parts[1]

				if !strings.HasPrefix(urlPath, basePath) {
					t.Errorf("buildNotificationURL() base path = %v, want %v", urlPath, basePath)
				}

				queryParts := strings.Split(queryParams, "&")
				for _, param := range queryParts {
					if !strings.Contains(urlPath, param) {
						t.Errorf("buildNotificationURL() missing query param %v in %v", param, urlPath)
					}
				}
			} else if urlPath != tt.wantPath {
				t.Errorf("buildNotificationURL() = %v, want %v", urlPath, tt.wantPath)
			}
		})
	}
}

// ExampleSend demonstrates how to use the Bark client in a real application.
// This is not an automated test, but rather a usage example and manual integration test.
func Example_send() {
	// Replace with your actual Bark key
	client, err := NewClient("", "YOUR-BARK-KEY")
	if err != nil {
		panic(err)
	}

	// Simple notification
	err = client.Send(context.Background(), "Hello from Go!")
	if err != nil {
		panic(err)
	}

	// Notification with newlines
	err = client.Send(context.Background(), "This is a multiline message.\nSecond line here.\nThird line here.",
		WithTitle("Multiline Test"))
	if err != nil {
		panic(err)
	}

	// Advanced notification
	err = client.Send(context.Background(), "Important meeting in 5 minutes!",
		WithTitle("Meeting Reminder"),
		WithSubtitle("Team Standup"),
		WithSound("bell"),
		WithTimeSensitive(),
	)
	if err != nil {
		panic(err)
	}
}

// TestIntegration_Send is a manual integration test.
// To run this test, set the BARK_KEY environment variable and use:
// go test -tags=integration -run TestIntegration_Send

func TestIntegration_Send(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	key := os.Getenv("BARK_KEY")
	if key == "" {
		t.Skip("BARK_KEY environment variable not set")
	}

	client, err := NewClient("", key)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		body string
		opts []Option
	}{
		{
			name: "simple notification",
			body: "Test message from integration test",
		},
		{
			name: "multiline notification",
			body: "Line 1\nLine 2\nLine 3",
			opts: []Option{
				WithTitle("Multiline Test"),
			},
		},
		{
			name: "full featured notification",
			body: "Test message with all features",
			opts: []Option{
				WithTitle("Integration Test"),
				WithSubtitle("All Features"),
				WithSound("bell"),
				WithTimeSensitive(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Send(context.Background(), tt.body, tt.opts...)
			if err != nil {
				t.Errorf("Send() error = %v", err)
			}
			// Add a delay between tests to avoid rate limiting
			time.Sleep(time.Second * 2)
		})
	}
}
