package zwift_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bzimmer/activity/zwift"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func newClient(t *testing.T, mux *http.ServeMux) (*zwift.Client, *httptest.Server) {
	a := assert.New(t)
	svr := httptest.NewServer(mux)

	endpoint := zwift.Endpoint()
	endpoint.AuthURL = svr.URL + "/auth"
	endpoint.TokenURL = svr.URL + "/token"

	client, err := zwift.NewClient(
		zwift.WithBaseURL(svr.URL),
		zwift.WithConfig(oauth2.Config{Endpoint: endpoint}),
		zwift.WithTokenCredentials("foo", "bar", time.Now().Add(time.Hour*24)),
		zwift.WithClientCredentials("what", "now?"))
	a.NoError(err)
	a.NotNil(client)
	return client, svr
}
