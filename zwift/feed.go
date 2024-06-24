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
	followeesFeedType feedType = "FOLLOWEES"
	justMeFeedType    feedType = "JUST_ME"
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

func (s *FeedService) feed(
	ctx context.Context, paginator *feedPaginator, spec activity.Pagination) ([]*Activity, error) {
	err := activity.Paginate(ctx, paginator, spec)
	if err != nil {
		return nil, err
	}
	return paginator.activities, nil
}

// JustMe returns all activities of the "just me" feed
func (s *FeedService) JustMe(
	ctx context.Context, includeInProgress bool, spec activity.Pagination) ([]*Activity, error) {
	return s.feed(ctx, &feedPaginator{
		service:           *s,
		includeInProgress: includeInProgress,
		feedType:          justMeFeedType,
		activities:        make([]*Activity, 0),
	}, spec)
}

// Favorites returns the feed for all riders marked as 'favorite'
func (s *FeedService) Favorites(
	ctx context.Context, includeInProgress bool, spec activity.Pagination) ([]*Activity, error) {
	return s.feed(ctx, &feedPaginator{
		service:           *s,
		includeInProgress: includeInProgress,
		feedType:          favoritesFeedType,
		activities:        make([]*Activity, 0),
	}, spec)
}

// Followees returns the feed for all riders which the logged in user follows
func (s *FeedService) Followees(
	ctx context.Context, includeInProgress bool, spec activity.Pagination) ([]*Activity, error) {
	return s.feed(ctx, &feedPaginator{
		service:           *s,
		includeInProgress: includeInProgress,
		feedType:          followeesFeedType,
		activities:        make([]*Activity, 0),
	}, spec)
}
