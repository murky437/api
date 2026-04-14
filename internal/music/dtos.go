package music

type SearchResponse struct {
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
}
