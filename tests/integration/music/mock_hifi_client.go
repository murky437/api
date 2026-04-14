package music

import (
	"murky_api/internal/music/hifiapi"
)

type MockHifiClient struct {
	SearchFunc func(params hifiapi.SearchParams) (*hifiapi.SearchResponse, error)
}

func (m *MockHifiClient) Search(params hifiapi.SearchParams) (*hifiapi.SearchResponse, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(params)
	}

	// Default mock response - return empty results
	return &hifiapi.SearchResponse{
		Data: hifiapi.SearchResponseData{
			Items: []hifiapi.SearchItem{},
		},
	}, nil
}
