package status

type GetResponse struct {
	Commit    string `json:"commit"`
	Timestamp string `json:"timestamp"`
}
