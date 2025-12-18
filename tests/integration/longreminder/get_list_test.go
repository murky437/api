package longreminder

import (
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/longreminder"
	"murky_api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetListUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/long-reminders", nil)
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

	req := httptest.NewRequest(http.MethodGet, "/long-reminders", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp longreminder.ListResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, 2, len(resp.Data))

	require.Equal(t, "longReminder1", resp.Data[0].Id)
	require.Equal(t, "Long reminder 1", resp.Data[0].Title)
	require.Equal(t, 40, resp.Data[0].IntervalDays)
	require.NotNil(t, resp.Data[0].CreatedAt)
	require.Equal(t, "2025-11-11T10:20:30Z", *resp.Data[0].UpdatedAt)
	require.Nil(t, resp.Data[0].LastRemindedAt)
	require.Nil(t, resp.Data[0].MarkedDoneAt)

	require.Equal(t, "longReminder2", resp.Data[1].Id)
	require.Equal(t, "Long reminder 2", resp.Data[1].Title)
	require.Equal(t, 40, resp.Data[1].IntervalDays)
	require.NotNil(t, resp.Data[1].CreatedAt)
	require.Equal(t, "2025-11-11T10:20:30Z", *resp.Data[1].UpdatedAt)
	require.Nil(t, resp.Data[1].LastRemindedAt)
	require.Nil(t, resp.Data[1].MarkedDoneAt)
}
