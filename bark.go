// Package gobark provides a Go SDK for the Bark push notification service.
// Bark is a simple and secure push notification tool that leverages APNs.
package gobark

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Client represents a Bark API client.
type Client struct {
	baseURL string
	key     string
	client  *http.Client
}

// NotificationLevel represents the level of notification importance.
type NotificationLevel string

const (
	// LevelActive is the default notification level.
	LevelActive NotificationLevel = "active"
	// LevelTimeSensitive represents time-sensitive notifications that can break through focus mode.
	LevelTimeSensitive NotificationLevel = "timeSensitive"
	// LevelPassive adds notification to the list without lighting up the screen.
	LevelPassive NotificationLevel = "passive"
	// LevelCritical represents critical alerts that ignore silent and do not disturb modes.
	LevelCritical NotificationLevel = "critical"

	defaultTitle = "无名消息"
)

// notification represents a Bark notification request.
type notification struct {
	title      string
	body       string
	subtitle   string
	icon       string
	sound      string
	level      NotificationLevel
	isCritical bool
}

// Option represents a function that modifies the notification request.
type Option func(*notification)

// NewClient creates a new Bark client with the specified base URL and key.
func NewClient(baseURL, key string) (*Client, error) {
	if baseURL == "" {
		baseURL = "https://api.day.app"
	}

	if key == "" {
		return nil, fmt.Errorf("bark key is required")
	}

	return &Client{
		baseURL: baseURL,
		key:     key,
		client:  &http.Client{},
	}, nil
}

// WithTitle sets the notification title.
func WithTitle(title string) Option {
	return func(n *notification) {
		n.title = title
	}
}

// WithSubtitle sets the notification subtitle.
func WithSubtitle(subtitle string) Option {
	return func(n *notification) {
		n.subtitle = subtitle
	}
}

// WithIcon sets the notification icon URL (iOS 15+ only).
func WithIcon(iconURL string) Option {
	return func(n *notification) {
		n.icon = iconURL
	}
}

// WithSound sets the notification sound.
func WithSound(sound string) Option {
	return func(n *notification) {
		n.sound = sound
	}
}

// WithTimeSensitive sets the notification as time-sensitive.
func WithTimeSensitive() Option {
	return func(n *notification) {
		n.level = LevelTimeSensitive
	}
}

// WithCriticalNotify sets the notification as a critical alert.
func WithCriticalNotify() Option {
	return func(n *notification) {
		n.level = LevelCritical
		n.isCritical = true
	}
}

// buildNotificationURL constructs the complete notification URL with all parameters
func (c *Client) buildNotificationURL(n *notification) string {
	// URL encode the body to handle special characters, especially newlines (\n)
	encodedBody := url.PathEscape(n.body)

	// Build the URL path based on available parameters
	urlPath := c.key
	if n.title != "" && n.subtitle != "" {
		urlPath = fmt.Sprintf("%s/%s/%s/%s", urlPath, url.PathEscape(n.title), url.PathEscape(n.subtitle), encodedBody)
	} else if n.title != "" {
		urlPath = fmt.Sprintf("%s/%s/%s", urlPath, url.PathEscape(n.title), encodedBody)
	} else {
		urlPath = fmt.Sprintf("%s/%s", urlPath, encodedBody)
	}

	// Build the query parameters for additional options
	query := url.Values{}
	if n.icon != "" {
		query.Set("icon", n.icon)
	}
	if n.sound != "" {
		query.Set("sound", n.sound)
	}
	if n.level != "" {
		query.Set("level", string(n.level))
	}
	if n.isCritical {
		query.Set("level", "critical")
	}

	// Construct the final URL
	apiURL := fmt.Sprintf("%s/%s", c.baseURL, urlPath)
	if len(query) > 0 {
		apiURL += "?" + query.Encode()
	}

	return apiURL
}

// Send sends a push notification through Bark.
// The body parameter is required and represents the main content of the notification.
// Additional options can be provided to customize the notification.
func (c *Client) Send(ctx context.Context, body string, opts ...Option) error {
	if body == "" {
		return fmt.Errorf("notification body is required")
	}

	n := &notification{
		title: defaultTitle,
		body:  body,
	}

	for _, opt := range opts {
		opt(n)
	}

	apiURL := c.buildNotificationURL(n)

	// Create and send the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
