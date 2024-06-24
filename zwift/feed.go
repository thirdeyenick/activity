package zwift

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bzimmer/activity"
)

const (
	feedPageCount = 30

	favoritesFeedType feedType = "FAVORITES"
)

type feedType string

// FeedService is the API for feed endpoints
type FeedService service

type feedPaginator struct {
	service           FeedService
	feedType          feedType
	includeInProgress bool
	activities        []*Activity
}

func (p *feedPaginator) PageSize() int {
	return feedPageCount
}

func (p *feedPaginator) Count() int {
	return len(p.activities)
}

func (p *feedPaginator) Do(ctx context.Context, spec activity.Pagination) (int, error) {
	// pagination uses the concept of page (based on strava), zwift uses an offset by row
	//  since pagination starts with page 1 (again, strava), subtract one from `start`
	count := spec.Count

	uri := fmt.Sprintf("api/activity-feed/feed/?limit=%d&includeInProgress=%t&feedType=%s", count, p.includeInProgress, p.feedType)
	req, err := p.service.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return 0, err
	}
	var acts []*Activity
	if err = p.service.client.do(req, &acts); err != nil {
		return 0, err
	}
	if spec.Total > 0 && len(p.activities)+len(acts) > spec.Total {
		acts = acts[:spec.Total-len(p.activities)]
	}
	p.activities = append(p.activities, acts...)
	return len(acts), nil
}

// Favorites returns the feed for all riders marked as 'favorite'
func (s *FeedService) Favorites(
	ctx context.Context, includeInProgress bool, spec activity.Pagination) ([]*Activity, error) {
	p := &feedPaginator{
		service:           *s,
		includeInProgress: includeInProgress,
		feedType:          favoritesFeedType,
		activities:        make([]*Activity, 0),
	}
	err := activity.Paginate(ctx, p, spec)
	if err != nil {
		return nil, err
	}
	return p.activities, nil
}
