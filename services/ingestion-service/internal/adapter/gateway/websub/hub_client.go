package websub

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

const (
	youtubeWebSubHubURL = "https://pubsubhubbub.appspot.com/subscribe"
	youtubeFeedBaseURL  = "https://www.youtube.com/xml/feeds/videos.xml"
)

// hubClient implements WebSubHub interface for YouTube
type hubClient struct {
	httpClient *http.Client
}

// NewHubClient creates a new WebSub hub client
func NewHubClient() gateway.WebSubHub {
	return &hubClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Subscribe subscribes to a YouTube channel's feed
func (c *hubClient) Subscribe(ctx context.Context, channelID valueobject.YouTubeChannelID, callbackURL string, leaseSeconds int) error {
	topicURL := c.buildTopicURL(channelID)
	
	params := url.Values{
		"hub.callback":      {callbackURL},
		"hub.topic":         {topicURL},
		"hub.verify":        {"async"},
		"hub.mode":          {"subscribe"},
		"hub.lease_seconds": {fmt.Sprintf("%d", leaseSeconds)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, youtubeWebSubHubURL, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send subscribe request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Unsubscribe unsubscribes from a YouTube channel's feed
func (c *hubClient) Unsubscribe(ctx context.Context, channelID valueobject.YouTubeChannelID, callbackURL string) error {
	topicURL := c.buildTopicURL(channelID)
	
	params := url.Values{
		"hub.callback": {callbackURL},
		"hub.topic":    {topicURL},
		"hub.verify":   {"async"},
		"hub.mode":     {"unsubscribe"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, youtubeWebSubHubURL, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send unsubscribe request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// buildTopicURL builds the YouTube channel feed URL
func (c *hubClient) buildTopicURL(channelID valueobject.YouTubeChannelID) string {
	return fmt.Sprintf("%s?channel_id=%s", youtubeFeedBaseURL, string(channelID))
}