package project

import (
	"bytes"
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/model"
	"murky_api/internal/project"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpdateNotesUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/projects/project-1/notes", nil)
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateNotesInvalidContentType(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-1/notes", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestUpdateNotesInvalidJson(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-1", strings.NewReader("{]"))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateNotesSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	body := project.UpdateNotesRequest{Notes: "asdasd"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-1/notes", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	p, err := model.GetProjectNotesByUserIdAndSlug(c.Db, 1, "project-1")
	require.NoError(t, err)
	require.Equal(t, body.Notes, p.Notes)
}

func TestUpdateOtherUserNotesNotFound(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	body := project.UpdateNotesRequest{Notes: "asd"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-4/notes", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}
