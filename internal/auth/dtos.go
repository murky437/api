package auth

import (
	"murky_api/internal/config"
	"murky_api/internal/validation"
	"net/mail"
	"slices"
	"strings"
)

type CreateTokensRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (request *CreateTokensRequest) Validate() *validation.Result {
	result := &validation.Result{
		GeneralErrors: []string{},
		FieldErrors:   make(map[string][]string),
	}

	if strings.TrimSpace(request.Username) == "" {
		result.FieldErrors["username"] = append(result.FieldErrors["username"], validation.NotBlankMessage)
	}

	if strings.TrimSpace(request.Password) == "" {
		result.FieldErrors["password"] = append(result.FieldErrors["password"], validation.NotBlankMessage)
	}

	if len(result.GeneralErrors) > 0 || len(result.FieldErrors) > 0 {
		return result
	}

	return nil
}

type CreateTokensResponse struct {
	AccessToken string `json:"accessToken"`
}

type RefreshAccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type SendGuestLoginLinkRequest struct {
	Email string `json:"email"`
	Url   string `json:"url"`
}

func (request *SendGuestLoginLinkRequest) Validate(conf *config.Config) *validation.Result {
	result := &validation.Result{
		GeneralErrors: []string{},
		FieldErrors:   make(map[string][]string),
	}

	if strings.TrimSpace(request.Email) == "" {
		result.FieldErrors["email"] = append(result.FieldErrors["email"], validation.NotBlankMessage)
	}

	if _, err := mail.ParseAddress(request.Email); err != nil {
		result.FieldErrors["email"] = append(result.FieldErrors["email"], "Needs to be a valid email address.")
	}

	if strings.TrimSpace(request.Url) == "" {
		result.FieldErrors["url"] = append(result.FieldErrors["url"], validation.NotBlankMessage)
	}

	if !slices.Contains(conf.AllowedOrigins, request.Url) {
		result.FieldErrors["url"] = append(result.FieldErrors["url"], "Needs to be a valid frontend url.")
	}

	if len(result.FieldErrors) > 0 || len(result.FieldErrors) > 0 {
		return result
	}

	return nil
}

type CreateTokenWithGuestTokenRequest struct {
	GuestToken string `json:"guestToken"`
}

type CreateTokenWithGuestTokenResponse struct {
	AccessToken string `json:"accessToken"`
}
