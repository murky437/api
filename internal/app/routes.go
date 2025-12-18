package app

import (
	"murky_api/internal/auth"
	"murky_api/internal/longreminder"
	"murky_api/internal/project"
	"murky_api/internal/routing"
	"murky_api/internal/status"
	"net/http"
)

func NewMux(c *Container) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/create-tokens", routing.Chain(
		auth.CreateTokens(c.Db, c.JwtService),
		routing.RequireJSON,
		routing.CorsMiddleware(c.Config)))
	mux.HandleFunc("POST /auth/refresh-access-token", routing.Chain(
		auth.RefreshAccessToken(c.Db, c.JwtService),
		routing.CorsMiddleware(c.Config)))
	mux.HandleFunc("POST /auth/delete-refresh-token", routing.Chain(
		auth.DeleteRefreshToken(c.Db),
		routing.CorsMiddleware(c.Config)))

	mux.HandleFunc("POST /auth/send-guest-login-link", routing.Chain(
		auth.SendGuestLoginLink(c.Db, c.GuestDb, c.Config, c.EmailClient),
		routing.RequireJSON,
		routing.CorsMiddleware(c.Config)))
	mux.HandleFunc("POST /auth/create-token-with-guest-token", routing.Chain(
		auth.CreateTokenWithGuestToken(c.GuestDb, c.JwtService),
		routing.RequireJSON,
		routing.CorsMiddleware(c.Config)))

	mux.HandleFunc("GET /status", routing.Chain(
		status.Get(c.Config),
		routing.CorsMiddleware(c.Config)))

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		routing.WriteNotFoundResponse(w)
	})

	protectedMux.HandleFunc("POST /projects", routing.Chain(project.Create(), routing.RequireJSON))
	protectedMux.HandleFunc("GET /projects", project.GetList())
	protectedMux.HandleFunc("GET /projects/{slug}", project.Get())
	protectedMux.HandleFunc("PUT /projects/{slug}", routing.Chain(project.Update(), routing.RequireJSON))
	protectedMux.HandleFunc("GET /projects/{slug}/notes", project.GetNotes())
	protectedMux.HandleFunc("PUT /projects/{slug}/notes", routing.Chain(project.UpdateNotes(), routing.RequireJSON))
	protectedMux.HandleFunc("DELETE /projects/{slug}", project.Delete())
	protectedMux.HandleFunc("PUT /projects/{slug}/sort-index", routing.Chain(project.UpdateSortIndex(), routing.RequireJSON))

	protectedMux.HandleFunc("POST /long-reminders", routing.Chain(longreminder.Create(), routing.RequireJSON))
	protectedMux.HandleFunc("GET /long-reminders", longreminder.GetList())
	protectedMux.HandleFunc("DELETE /long-reminders/{id}", longreminder.Delete())
	protectedMux.HandleFunc("PUT /long-reminders/{id}", routing.Chain(longreminder.Update(), routing.RequireJSON))

	mux.HandleFunc("/", routing.Chain(
		protectedMux.ServeHTTP,
		routing.RequireAuth(c.Db, c.JwtService, c.GuestDb),
		routing.CorsMiddleware(c.Config),
	))

	return mux
}
