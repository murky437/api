package longreminder

import (
	"bytes"
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/constants"
	"murky_api/internal/longreminder"
	"murky_api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpdateUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/long-reminders/longReminder1", nil)
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateTitleAndIntervalDaysSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	title := "newTitle"
	intervalDays := 5
	body := longreminder.UpdateRequest{Title: &title, IntervalDays: &intervalDays}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/long-reminders/longReminder1", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	r, err := model.GetLongReminderByUserIdAndId(c.Db, 1, "longReminder1")
	require.NoError(t, err)
	require.Equal(t, r.Title, title)
	require.Equal(t, r.IntervalDays, intervalDays)
	require.Nil(t, r.MarkedDoneAt)
	require.Equal(t, r.IsEnabled, true)
}

func TestUpdateMarkedDoneAtSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	markedDoneAt := "2025-11-11T01:02:03.123Z"
	body := longreminder.UpdateRequest{MarkedDoneAt: &markedDoneAt}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/long-reminders/longReminder1", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	r, err := model.GetLongReminderByUserIdAndId(c.Db, 1, "longReminder1")
	require.NoError(t, err)
	require.Equal(t, r.Title, "Long reminder 1")
	require.Equal(t, r.IntervalDays, 40)
	markedDoneAtTime, err := time.Parse(constants.ApiDateTimeFormat, markedDoneAt)
	require.NoError(t, err)
	require.Equal(t, r.MarkedDoneAt, &markedDoneAtTime)
	require.Equal(t, r.IsEnabled, true)
}

func TestUpdateIsEnabledSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	isEnabled := false
	body := longreminder.UpdateRequest{IsEnabled: &isEnabled}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/long-reminders/longReminder1", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	r, err := model.GetLongReminderByUserIdAndId(c.Db, 1, "longReminder1")
	require.NoError(t, err)
	require.Equal(t, r.Title, "Long reminder 1")
	require.Equal(t, r.IntervalDays, 40)
	require.Nil(t, r.MarkedDoneAt)
	require.Equal(t, r.IsEnabled, false)
}

// TODO: more tests
