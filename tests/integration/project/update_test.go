package project

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/model"
	"murky_api/internal/project"
	"murky_api/internal/routing"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpdateUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/projects/project-1", nil)
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateInvalidContentType(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-1", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestUpdateInvalidJson(t *testing.T) {
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

func TestUpdateValidationError(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-1", strings.NewReader("{}"))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp routing.ValidationErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []string{
		"Must not be blank.",
		"Has to be at least 2 characters long.",
		"Must only contain letters, numbers, dashes and spaces.",
	}, resp.FieldErrors["title"])
	require.Equal(t, []string{
		"Must not be blank.",
		"Has to be at least 2 characters long.",
		"Must only contain letters, numbers and dashes.",
	}, resp.FieldErrors["slug"])
}

func TestUpdateNotUniqueSlugError(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	body := project.UpdateRequest{Title: "Test Project", Slug: "project-2"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-1", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp routing.ValidationErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []string{"A project with this slug already exists. Slug must be unique."}, resp.FieldErrors["slug"])
}

func TestUpdateSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	body := project.UpdateRequest{Title: "Test Project", Slug: "test-project"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-1", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	p, err := model.GetProjectByUserIdAndSlug(c.Db, 1, "project-1")
	require.ErrorIs(t, err, sql.ErrNoRows)
	p, err = model.GetProjectByUserIdAndSlug(c.Db, 1, "test-project")
	require.NoError(t, err)
	require.Equal(t, p.Title, "Test Project")
	require.Equal(t, p.Slug, "test-project")
}

func TestUpdateOtherUserNotFound(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	body := project.UpdateRequest{Title: "Test Project", Slug: "test-project"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/projects/project-4", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}
