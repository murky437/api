package project

import (
	"database/sql"
	"murky_api/internal/model"
	"murky_api/internal/validation"
	"regexp"
	"strings"
)

type CreateRequest struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

func (request *CreateRequest) Validate(db *sql.DB) (*validation.Result, error) {
	result := &validation.Result{
		GeneralErrors: []string{},
		FieldErrors:   make(map[string][]string),
	}

	if strings.TrimSpace(request.Title) == "" {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], validation.NotBlankMessage)
	}

	if strings.TrimSpace(request.Slug) == "" {
		result.FieldErrors["slug"] = append(result.FieldErrors["slug"], validation.NotBlankMessage)
	}

	if len(request.Title) < 2 {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], "Has to be at least 2 characters long.")
	}

	if len(request.Slug) < 2 {
		result.FieldErrors["slug"] = append(result.FieldErrors["slug"], "Has to be at least 2 characters long.")
	}

	matchedTitle := regexp.MustCompile(`^[a-zA-Z0-9\- ]+$`).MatchString(request.Title)
	if !matchedTitle {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], "Must only contain letters, numbers, dashes and spaces.")
	}

	matchedSlug := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString(request.Slug)
	if !matchedSlug {
		result.FieldErrors["slug"] = append(result.FieldErrors["slug"], "Must only contain letters, numbers and dashes.")
	}

	if len(result.GeneralErrors) > 0 || len(result.FieldErrors) > 0 {
		return result, nil
	}

	isUnique, err := model.IsProjectSlugUnique(db, request.Slug, model.SlugCheckOptions{})
	if err != nil {
		return nil, err
	}
	if !isUnique {
		result.FieldErrors["slug"] = append(
			result.FieldErrors["slug"],
			"A project with this slug already exists. Slug must be unique.",
		)
	}

	if len(result.GeneralErrors) > 0 || len(result.FieldErrors) > 0 {
		return result, nil
	}

	return nil, nil
}

type CreateResponse struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type ListResponse struct {
	Data []ListItem `json:"data"`
}

type ListItem struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type GetResponse struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type GetNotesResponse struct {
	Notes string `json:"notes"`
}

type UpdateRequest struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

func (request *UpdateRequest) Validate(db *sql.DB, currentSlug string) (*validation.Result, error) {
	result := &validation.Result{
		GeneralErrors: []string{},
		FieldErrors:   make(map[string][]string),
	}

	if strings.TrimSpace(request.Title) == "" {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], validation.NotBlankMessage)
	}

	if strings.TrimSpace(request.Slug) == "" {
		result.FieldErrors["slug"] = append(result.FieldErrors["slug"], validation.NotBlankMessage)
	}

	if len(request.Title) < 2 {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], "Has to be at least 2 characters long.")
	}

	if len(request.Slug) < 2 {
		result.FieldErrors["slug"] = append(result.FieldErrors["slug"], "Has to be at least 2 characters long.")
	}

	matchedTitle := regexp.MustCompile(`^[a-zA-Z0-9\- ]+$`).MatchString(request.Title)
	if !matchedTitle {
		result.FieldErrors["title"] = append(result.FieldErrors["title"], "Must only contain letters, numbers, dashes and spaces.")
	}

	matchedSlug := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString(request.Slug)
	if !matchedSlug {
		result.FieldErrors["slug"] = append(result.FieldErrors["slug"], "Must only contain letters, numbers and dashes.")
	}

	if len(result.GeneralErrors) > 0 || len(result.FieldErrors) > 0 {
		return result, nil
	}

	isUnique, err := model.IsProjectSlugUnique(db, request.Slug, model.SlugCheckOptions{CurrentSlug: currentSlug})
	if err != nil {
		return nil, err
	}
	if !isUnique {
		result.FieldErrors["slug"] = append(
			result.FieldErrors["slug"],
			"A project with this slug already exists. Slug must be unique.",
		)
	}

	if len(result.GeneralErrors) > 0 || len(result.FieldErrors) > 0 {
		return result, nil
	}

	return nil, nil
}

type UpdateNotesRequest struct {
	Notes string `json:"notes"`
}

type UpdateSortIndexRequest struct {
	SortIndex int `json:"sortIndex"`
}

func (request *UpdateSortIndexRequest) Validate(db *sql.DB, userId int) (*validation.Result, error) {
	result := &validation.Result{
		GeneralErrors: []string{},
		FieldErrors:   make(map[string][]string),
	}

	projectCount, err := model.GetProjectCount(db, userId)
	if err != nil {
		return nil, err
	}

	if request.SortIndex < 0 || request.SortIndex >= projectCount {
		result.FieldErrors["sortIndex"] = append(result.FieldErrors["sortIndex"], "Invalid sort index.")
	}

	if len(result.GeneralErrors) > 0 || len(result.FieldErrors) > 0 {
		return result, nil
	}

	return nil, nil
}
