package music

import (
	"murky_api/internal/music/hifiapi"
	"murky_api/internal/routing"
	"murky_api/internal/validation"
	"net/http"
	"net/url"
)

func GetSearch(hifiClient hifiapi.HifiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		track := query.Get("track")
		artist := query.Get("artist")

		if track == "" && artist == "" {
			validationResult := validation.Result{
				GeneralErrors: []string{"At least one of 'track' or 'artist' query parameters must be provided"},
			}
			routing.WriteValidationErrorResponse(w, validationResult)
			return
		}

		params := hifiapi.SearchParams{
			Track:  track,
			Artist: artist,
		}

		searchResp, err := hifiClient.Search(params)
		if err != nil {
			if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
				routing.WriteJsonResponse(w, http.StatusGatewayTimeout, routing.GeneralErrorResponse{Message: "Hifi API request timed out"})
				return
			}
			routing.WriteJsonResponse(w, http.StatusBadGateway, routing.GeneralErrorResponse{Message: "Failed to fetch data from Hifi API"})
			return
		}

		resp := SearchResponse{
			Items: make([]SearchItem, len(searchResp.Data.Items)),
		}

		for i, item := range searchResp.Data.Items {
			resp.Items[i] = SearchItem{
				ID:       item.ID,
				Title:    item.Title,
				Duration: item.Duration,
				Artist:   item.Artist.Name,
				Album:    item.Album.Title,
			}
		}

		routing.WriteJsonResponse(w, http.StatusOK, resp)
	}
}
