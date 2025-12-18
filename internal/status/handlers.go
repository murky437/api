package status

import (
	"encoding/json"
	"log"
	"murky_api/internal/config"
	"murky_api/internal/routing"
	"net/http"
	"os"
)

func Get(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(conf.DeployFilePath)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		var resp GetResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		routing.WriteJsonResponse(w, http.StatusOK, resp)
	}
}
