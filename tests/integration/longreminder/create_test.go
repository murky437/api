package longreminder

import (
	"bytes"
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/longreminder"
	"murky_api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	reminderId := uuid.New().String()

	body := longreminder.CreateRequest{
		Id:           reminderId,
		Title:        "Test LongReminder",
		IntervalDays: 40,
	}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/long-reminders", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	r, err := model.GetLongReminderByUserIdAndId(c.Db, 1, reminderId)
	require.NoError(t, err)
	require.Equal(t, r.Title, "Test LongReminder")
	require.Equal(t, r.IntervalDays, 40)
}

// TODO: more tests
