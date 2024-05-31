package cyclinganalytics_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/activity"
	"github.com/bzimmer/activity/cyclinganalytics"
)

func TestUpload(t *testing.T) {
	tests := []struct {
		name, err string
		user      cyclinganalytics.UserID
		file      *activity.File
	}{
		{
			name: "uploader error",
			err:  "missing upload file",
		},
		{
			name: "uploader success me",
			file: &activity.File{
				Filename: "/path/to/foo.gpx",
				Name:     "foo.gpx",
				Format:   activity.FormatGPX,
				Reader:   bytes.NewBufferString("<gpx></gpx>"),
			},
		},
		{
			name: "uploader success user",
			user: 2298801,
			file: &activity.File{
				Filename: "/path/to/foo.gpx",
				Name:     "foo.gpx",
				Format:   activity.FormatGPX,
				Reader:   bytes.NewBufferString("<gpx></gpx>"),
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			mux := http.NewServeMux()
			mux.HandleFunc("/me/upload", func(w http.ResponseWriter, _ *http.Request) {
				enc := json.NewEncoder(w)
				a.NoError(enc.Encode(&cyclinganalytics.Upload{}))
			})
			mux.HandleFunc("/user/2298801/upload", func(w http.ResponseWriter, _ *http.Request) {
				enc := json.NewEncoder(w)
				a.NoError(enc.Encode(&cyclinganalytics.Upload{}))
			})

			svr := httptest.NewServer(mux)
			defer svr.Close()

			client, err := cyclinganalytics.NewClient(
				cyclinganalytics.WithBaseURL(svr.URL),
				cyclinganalytics.WithTokenCredentials("fooKey", "barToken", time.Time{}))
			a.NoError(err)
			uploader := client.Uploader()
			a.NotNil(uploader)
			var upload activity.Upload
			if tt.user != cyclinganalytics.Me {
				upload, err = client.Rides.UploadWithUser(context.Background(), tt.user, tt.file)
			} else {
				upload, err = uploader.Upload(context.Background(), tt.file)
			}
			if tt.err != "" {
				a.Error(err)
				a.Nil(upload)
				a.Contains(err.Error(), tt.err)
				return
			}
			a.NoError(err)
			a.NotNil(upload)
		})
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name, err string
		user      cyclinganalytics.UserID
		upload    activity.UploadID
	}{
		{
			name:   "status success me",
			user:   cyclinganalytics.Me,
			upload: 1891982,
		},
		{
			name:   "status success me",
			user:   882722,
			upload: 1891982,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			mux := http.NewServeMux()
			mux.HandleFunc("/me/upload/1891982", func(w http.ResponseWriter, _ *http.Request) {
				enc := json.NewEncoder(w)
				a.NoError(enc.Encode(&cyclinganalytics.Upload{}))
			})
			mux.HandleFunc("/user/882722/upload/1891982", func(w http.ResponseWriter, _ *http.Request) {
				enc := json.NewEncoder(w)
				a.NoError(enc.Encode(&cyclinganalytics.Upload{}))
			})

			svr := httptest.NewServer(mux)
			defer svr.Close()

			client, err := cyclinganalytics.NewClient(
				cyclinganalytics.WithBaseURL(svr.URL),
				cyclinganalytics.WithTokenCredentials("fooKey", "barToken", time.Time{}))
			a.NoError(err)
			uploader := client.Uploader()
			a.NotNil(uploader)
			var upload activity.Upload
			if tt.user == cyclinganalytics.Me {
				upload, err = uploader.Status(context.Background(), tt.upload)
			} else {
				upload, err = client.Rides.StatusWithUser(context.Background(), tt.user, int64(tt.upload))
			}
			if tt.err != "" {
				a.Error(err)
				a.Nil(upload)
				a.Contains(err.Error(), tt.err)
				return
			}
			a.NoError(err)
			a.NotNil(upload)
		})
	}
}

func TestUploadIdentifier(t *testing.T) {
	a := assert.New(t)
	u := &cyclinganalytics.Upload{ID: 1992822, Status: "processing"}
	a.Equal(activity.UploadID(1992822), u.Identifier())
	a.False(u.Done())
	u.Status = "done"
	a.True(u.Done())
}
