package longreminder

import (
	"database/sql"
	"fmt"
	"murky_api/internal/constants"
	"murky_api/internal/validation"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	IntervalDays int    `json:"intervalDays"`
}

func (request *CreateRequest) Validate(db *sql.DB) *validation.Result {
	result := &validation.Result{
		GeneralErrors: []string{},
		FieldErrors:   make(map[string][]string),
	}

	if strings.TrimSpace(request.Id) == "" {
		result.FieldErrors["id"] = append(result.FieldErrors["id"], validation.NotBlankMessage)
	}

	if strings.TrimSpace(request.Title) == "" {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], validation.NotBlankMessage)
	}

	if _, err := uuid.Parse(request.Id); err != nil {
		result.FieldErrors["id"] = append(result.FieldErrors["id"], "Invalid uuid.")
	}

	if len(request.Title) < 2 {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], "Has to be at least 2 characters long.")
	}

	matchedTitle := regexp.MustCompile(`^[a-zA-Z0-9\- ]+$`).MatchString(request.Title)
	if !matchedTitle {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], "Must only contain letters, numbers, dashes and spaces.")
	}

	if request.IntervalDays < 1 {
		result.FieldErrors["intervalDays"] = append(result.FieldErrors["intervalDays"], "Has to be at least 1.")
	}

	if len(result.GeneralErrors) > 0 || len(result.FieldErrors) > 0 {
		return result
	}

	return nil
}

type ListResponse struct {
	Data []ListItem `json:"data"`
}

type ListItem struct {
	Id             string  `json:"id"`
	Title          string  `json:"title"`
	IntervalDays   int     `json:"intervalDays"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      *string `json:"updatedAt"`
	LastRemindedAt *string `json:"lastRemindedAt"`
	MarkedDoneAt   *string `json:"markedDoneAt"`
	IsEnabled      bool    `json:"isEnabled"`
}
type UpdateRequest struct {
	Title        *string `json:"title"`
	IntervalDays *int    `json:"intervalDays"`
	MarkedDoneAt *string `json:"markedDoneAt"`
	IsEnabled    *bool   `json:"isEnabled"`
}

func (request *UpdateRequest) Validate(db *sql.DB) *validation.Result {
	result := &validation.Result{
		GeneralErrors: []string{},
		FieldErrors:   make(map[string][]string),
	}

	if request.Title != nil {
		if strings.TrimSpace(*request.Title) == "" {
			result.FieldErrors["title"] = append(result.FieldErrors["title"], validation.NotBlankMessage)
		}

		if len(*request.Title) < 2 {
			result.FieldErrors["title"] = append(result.FieldErrors["title"], "Has to be at least 2 characters long.")
		}

		matchedTitle := regexp.MustCompile(`^[a-zA-Z0-9\- ]+$`).MatchString(*request.Title)
		if !matchedTitle {
			result.FieldErrors["title"] = append(result.FieldErrors["title"], "Must only contain letters, numbers, dashes and spaces.")
		}
	}

	if request.IntervalDays != nil {
		if *request.IntervalDays < 1 {
			result.FieldErrors["intervalDays"] = append(result.FieldErrors["intervalDays"], "Has to be at least 1.")
		}
	}

	if request.MarkedDoneAt != nil {
		if _, err := time.Parse(constants.ApiDateTimeFormat, *request.MarkedDoneAt); err != nil {
			result.FieldErrors["markedDoneAt"] = append(
				result.FieldErrors["markedDoneAt"],
				fmt.Sprintf("Date needs to be in format %s.", constants.ApiDateTimeFormat),
			)
		}
	}

	if len(result.GeneralErrors) > 0 || len(result.FieldErrors) > 0 {
		return result
	}

	return nil
}
