package music

import (
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/model"
	"murky_api/internal/music"
	"murky_api/internal/music/hifiapi"
	"murky_api/internal/routing"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetSearchUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/music/search?track=test", nil)
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetSearchMissingParameters(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/music/search", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp routing.ValidationErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []string{"'track' query parameter is required"}, resp.GeneralErrors)
}

func TestGetSearchWithTrackParameter(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	// Set mock client on container
	c.HifiClient = &MockHifiClient{
		SearchFunc: func(params hifiapi.SearchParams) (*hifiapi.SearchResponse, error) {
			// Verify correct parameters are passed
			require.Equal(t, "test", params.Track)

			// Return test data
			return &hifiapi.SearchResponse{
				Data: hifiapi.SearchResponseData{
					Items: []hifiapi.SearchItem{
						{
							ID:       123,
							Title:    "Test Track",
							Duration: 180,
							Artist:   hifiapi.Artist{Name: "Test Artist"},
							Album:    hifiapi.Album{Title: "Test Album"},
						},
					},
				},
			}, nil
		},
	}

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/music/search?track=test", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp music.SearchResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	// Verify response matches our mock data
	require.Len(t, resp.Items, 1)
	require.Equal(t, 123, resp.Items[0].ID)
	require.Equal(t, "Test Track", resp.Items[0].Title)
	require.Equal(t, "Test Artist", resp.Items[0].Artist)
	require.Equal(t, "Test Album", resp.Items[0].Album)
}
