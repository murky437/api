package project

import (
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/status"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	c.Config.DeployFilePath = "testdata/deploy.json"
	defer c.Close()

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp status.GetResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, status.GetResponse{
		Commit:    "testhsh",
		Timestamp: "20251024T151821",
	}, resp)
}

// TODO: more tests
