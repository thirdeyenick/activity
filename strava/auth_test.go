package strava_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client, err := newClient(http.StatusOK, "refresh.json")
	a.NoError(err)

	ctx := context.Background()
	tokens, err := client.Auth.Refresh(ctx)
	a.NoError(err, "failed to refresh")
	a.NotNil(tokens)
}
