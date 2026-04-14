package hifiapi

type SearchParams struct {
	Track    string `json:"s,omitempty"`
	Artist   string `json:"a,omitempty"`
	Album    string `json:"al,omitempty"`
	Video    string `json:"v,omitempty"`
	Playlist string `json:"p,omitempty"`
	ISRC     string `json:"i,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type SearchResponseData struct {
	Limit              int          `json:"limit"`
	Offset             int          `json:"offset"`
	TotalNumberOfItems int          `json:"totalNumberOfItems"`
	Items              []SearchItem `json:"items"`
}

type SearchResponse struct {
	Version string             `json:"version"`
	Data    SearchResponseData `json:"data"`
}

type SearchItem struct {
	ID                        int      `json:"id"`
	Title                     string   `json:"title"`
	Duration                  int      `json:"duration"`
	ReplayGain                float64  `json:"replayGain"`
	Peak                      float64  `json:"peak"`
	AllowStreaming            bool     `json:"allowStreaming"`
	StreamReady               bool     `json:"streamReady"`
	PayToStream               bool     `json:"payToStream"`
	AdSupportedStreamReady    bool     `json:"adSupportedStreamReady"`
	DjReady                   bool     `json:"djReady"`
	StemReady                 bool     `json:"stemReady"`
	StreamStartDate           string   `json:"streamStartDate"`
	PremiumStreamingOnly      bool     `json:"premiumStreamingOnly"`
	TrackNumber               int      `json:"trackNumber"`
	VolumeNumber              int      `json:"volumeNumber"`
	Version                   *string  `json:"version"`
	Popularity                int      `json:"popularity"`
	Copyright                 string   `json:"copyright"`
	ReleaseDate               string   `json:"releaseDate"`
	Url                       string   `json:"url"`
	Isrc                      string   `json:"isrc"`
	Editable                  bool     `json:"editable"`
	Explicit                  bool     `json:"explicit"`
	Type                      string   `json:"type"`
	BitDepth                  int      `json:"bitDepth"`
	SampleRate                float64  `json:"sampleRate"`
	AllowStreamingResolutions []string `json:"allowStreamingResolutions"`
	AudioModes                []string `json:"audioModes"`
	AudioQuality              string   `json:"audioQuality"`
	Artist                    Artist   `json:"artist"`
	Album                     Album    `json:"album"`
}

type Artist struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Picture string `json:"picture"`
	Url     string `json:"url"`
}

type Album struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"releaseDate"`
	Cover       string `json:"cover"`
	Url         string `json:"url"`
}
