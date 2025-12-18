package routing

import (
	"encoding/json"
	"fmt"
	"log"
	"murky_api/internal/validation"
	"net/http"
)

const (
	NotFoundMessage                        = "Not found."
	ContentTypeNotJsonMessage              = "Content-Type must be application/json."
	InvalidJsonMessage                     = "Invalid JSON."
	InternalServerErrorMessage             = "Internal server error."
	UnauthorizedMessage                    = "Unauthorized."
	TooManyRequestsMessage                 = "Too Many Requests."
	TooManyRequestsMessageWithRetrySeconds = "Too Many Requests. Try again in %d seconds."
)

type GeneralErrorResponse struct {
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	GeneralErrors []string            `json:"generalErrors"`
	FieldErrors   map[string][]string `json:"fieldErrors"`
}

func WriteJsonResponse(w http.ResponseWriter, status int, content any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if content != nil {
		err := json.NewEncoder(w).Encode(content)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func WriteNotFoundResponse(w http.ResponseWriter) {
	WriteJsonResponse(w, http.StatusNotFound, GeneralErrorResponse{Message: NotFoundMessage})
}

func WriteInternalServerErrorResponse(w http.ResponseWriter) {
	WriteJsonResponse(w, http.StatusInternalServerError, GeneralErrorResponse{Message: InternalServerErrorMessage})
}

func WriteContentTypeNotJsonResponse(w http.ResponseWriter) {
	WriteJsonResponse(w, http.StatusUnsupportedMediaType, GeneralErrorResponse{Message: ContentTypeNotJsonMessage})
}

func WriteInvalidJsonResponse(w http.ResponseWriter) {
	WriteJsonResponse(w, http.StatusBadRequest, GeneralErrorResponse{Message: InvalidJsonMessage})
}

func WriteValidationErrorResponse(w http.ResponseWriter, validationResult validation.Result) {
	WriteJsonResponse(w, http.StatusUnprocessableEntity, ValidationErrorResponse(validationResult))
}

func WriteUnauthorizedResponse(w http.ResponseWriter) {
	WriteJsonResponse(w, http.StatusUnauthorized, GeneralErrorResponse{Message: UnauthorizedMessage})
}

func WriteTooManyRequestsResponse(w http.ResponseWriter, retryAfterSeconds int) {
	if retryAfterSeconds > 0 {
		w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfterSeconds))
		WriteJsonResponse(w, http.StatusTooManyRequests, GeneralErrorResponse{Message: fmt.Sprintf(TooManyRequestsMessageWithRetrySeconds, retryAfterSeconds)})
		return
	}
	WriteJsonResponse(w, http.StatusTooManyRequests, GeneralErrorResponse{Message: TooManyRequestsMessage})
}

func WriteGeneralTooManyRequestsResponse(w http.ResponseWriter, message string) {
	WriteJsonResponse(w, http.StatusTooManyRequests, GeneralErrorResponse{Message: message})
}

func WriteGeneralErrorResponse(w http.ResponseWriter, message string) {
	WriteJsonResponse(w, http.StatusBadRequest, GeneralErrorResponse{Message: message})
}
