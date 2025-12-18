package project

import (
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/model"
	"murky_api/internal/project"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetListUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetListSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp project.ListResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []project.ListItem{
		project.ListItem{Title: "Project 1", Slug: "project-1"},
		project.ListItem{Title: "Project 2", Slug: "project-2"},
		project.ListItem{Title: "Project 3", Slug: "project-3"},
	}, resp.Data)
}
