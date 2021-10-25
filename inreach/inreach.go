package inreach

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

//go:generate genwith --client --do --decoder xml --package inreach

/*
 Information about the InReach KML feed is available at:
  https://support.garmin.com/en-US/?faq=tdlDCyo1fJ5UxjUbA9rMY8
*/

const (
	_baseURL   = "https://share.garmin.com"
	DateFormat = "2006-01-02T00:00z"
)

type Client struct {
	client  *http.Client
	baseURL string

	Feed *FeedService
}

// APIOption for configuring API requests
type APIOption func(url.Values) error

// WithBaseURL specifies the base url
func WithBaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// WithDateRange sets the date range for events
func WithDateRange(before, after time.Time) APIOption {
	return func(v url.Values) error {
		if !before.IsZero() && !after.IsZero() {
			if after.After(before) {
				return errors.New("invalid date range")
			}
		}
		if !before.IsZero() {
			v.Set("d2", before.Format(DateFormat))
		}
		if !after.IsZero() {
			v.Set("d1", after.Format(DateFormat))
		}
		return nil
	}
}

func withServices() Option {
	return func(c *Client) error {
		c.Feed = &FeedService{c}
		if c.baseURL == "" {
			c.baseURL = _baseURL
		}
		return nil
	}
}

func (c *Client) newRequest(ctx context.Context, method, uri string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", c.baseURL, uri))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
