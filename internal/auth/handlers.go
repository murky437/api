package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"murky_api/internal/config"
	"murky_api/internal/constants"
	"murky_api/internal/email"
	"murky_api/internal/jwt"
	"murky_api/internal/model"
	"murky_api/internal/routing"
	"murky_api/internal/validation"
	"net/http"
	"time"
)

func CreateTokens(db *sql.DB, jwtService jwt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTokensRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			routing.WriteInvalidJsonResponse(w)
			return
		}

		validationResult := req.Validate()
		if validationResult != nil {
			routing.WriteValidationErrorResponse(w, *validationResult)
			return
		}

		user, err := model.FindUserByUsername(db, req.Username)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		if user == nil || !isPasswordValid(user.Password, req.Password) {
			routing.WriteValidationErrorResponse(w, validation.Result{GeneralErrors: []string{"Invalid credentials"}})
			return
		}

		accessToken, err := jwtService.CreateAccessToken(*user, time.Now().UTC().Add(constants.AccessTokenDuration))
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}
		expiresAt := time.Now().UTC().Add(constants.RefreshTokenDuration)

		refreshTokenString, err := jwtService.CreateRefreshToken(req.Username, expiresAt)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		refreshToken := model.RefreshToken{
			UserId:    user.Id,
			Jwt:       refreshTokenString,
			ExpiresAt: expiresAt.Format(constants.SqliteDateTimeFormat),
		}

		if err := model.SaveRefreshToken(db, refreshToken); err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshTokenString,
			Path:     "/", // TODO: set a different path, so that its not sent for every request?
			HttpOnly: true,
			Secure:   false, // TODO: set to secure based on env (dev = false, prod = true)
			SameSite: http.SameSiteLaxMode,
			Expires:  expiresAt,
		})

		resp := CreateTokensResponse{
			AccessToken: accessToken,
		}

		routing.WriteJsonResponse(w, http.StatusOK, resp)
	}
}

func RefreshAccessToken(db *sql.DB, jwtService jwt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")

		if err != nil || cookie.Value == "" {
			routing.WriteUnauthorizedResponse(w)
			return
		}

		claims, err := jwtService.ParseRefreshToken(cookie.Value)
		if err != nil {
			routing.WriteUnauthorizedResponse(w)
			return
		}

		user, err := model.FindUserByUsername(db, claims.Username)
		if err != nil || user == nil {
			routing.WriteUnauthorizedResponse(w)
			return
		}

		dbToken, err := model.FindRefreshTokenByJwt(db, cookie.Value)
		if err != nil || dbToken == nil {
			routing.WriteUnauthorizedResponse(w)
			return
		}

		expiresAt, err := time.Parse(constants.SqliteDateTimeFormat, dbToken.ExpiresAt)
		if err != nil || time.Now().UTC().After(expiresAt) {
			routing.WriteUnauthorizedResponse(w)
			return
		}

		newAccessToken, err := jwtService.CreateAccessToken(*user, time.Now().UTC().Add(constants.AccessTokenDuration))
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		resp := RefreshAccessTokenResponse{
			AccessToken: newAccessToken,
		}

		routing.WriteJsonResponse(w, http.StatusOK, resp)
	}
}

func DeleteRefreshToken(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")

		if err != nil || cookie.Value == "" {
			routing.WriteJsonResponse(w, http.StatusNoContent, nil)
			return
		}

		err = model.DeleteRefreshToken(db, cookie.Value)
		if err != nil {
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1})

		routing.WriteJsonResponse(w, http.StatusNoContent, nil)
	}
}

func SendGuestLoginLink(mainDb *sql.DB, guestDb *sql.DB, conf *config.Config, emailClient email.Client) http.HandlerFunc {
	limiter := NewGuestLoginLimiter()
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendGuestLoginLinkRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			routing.WriteInvalidJsonResponse(w)
			return
		}

		validationResult := req.Validate(conf)
		if validationResult != nil {
			routing.WriteValidationErrorResponse(w, *validationResult)
			return
		}

		dailyBrevoSends, err := model.GetDailyBrevoSendsNumber(mainDb)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		if dailyBrevoSends >= 300 {
			routing.WriteGeneralTooManyRequestsResponse(w, "Guest logins are unavailable for the rest of today. Please try again tomorrow.")
			return
		}

		allowed, retryAfterSeconds := limiter.Reserve(req.Email, r)
		if !allowed {
			routing.WriteTooManyRequestsResponse(w, retryAfterSeconds)
			return
		}

		guestLoginToken, err := model.GetOrCreateGuestLoginToken(guestDb, req.Email)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		err = emailClient.SendWithTemplate(
			guestLoginToken.Email,
			fmt.Sprintf("Guest login link to %s", req.Url),
			1,
			map[string]any{
				"link": fmt.Sprintf("%s/token-login/%s", req.Url, guestLoginToken.Token),
			},
		)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		err = model.LogDailyBrevoSend(mainDb)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		routing.WriteJsonResponse(w, http.StatusNoContent, nil)
	}
}

func CreateTokenWithGuestToken(db *sql.DB, jwtService jwt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTokenWithGuestTokenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			routing.WriteInvalidJsonResponse(w)
			return
		}

		if req.GuestToken == "" {
			routing.WriteGeneralErrorResponse(w, "Invalid token.")
			return
		}

		guestLoginToken, err := model.GetGuestLoginTokenByToken(db, req.GuestToken)
		if err != nil || guestLoginToken == nil {
			routing.WriteGeneralErrorResponse(w, "Token has expired.")
			return
		}

		user, err := model.FindUserByUsername(db, guestLoginToken.Email)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		if user == nil {
			user, err = model.CreateGuestUser(db, guestLoginToken.Email)
			if err != nil {
				log.Println(err)
				routing.WriteInternalServerErrorResponse(w)
				return
			}
		}

		accessToken, err := jwtService.CreateAccessToken(*user, time.Now().UTC().Add(constants.AccessTokenDuration))
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		err = model.DeleteGuestLoginTokenByToken(db, guestLoginToken.Token)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		resp := CreateTokenWithGuestTokenResponse{
			AccessToken: accessToken,
		}

		routing.WriteJsonResponse(w, http.StatusOK, resp)
	}
}
