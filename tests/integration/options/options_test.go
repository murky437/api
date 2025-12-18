package options

import (
	"murky_api/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	c.Config.AllowedOrigins = []string{"http://localhost:5173", "http://localhost:8081"}
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
	require.Equal(t, "http://localhost:5173", rr.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
}
