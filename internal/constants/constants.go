package constants

import "time"

const (
	SqliteDateTimeFormat = "2006-01-02T15:04:05.000Z"
	ApiDateTimeFormat    = "2006-01-02T15:04:05Z"
	AccessTokenDuration  = time.Minute * 15
	RefreshTokenDuration = time.Hour * 24 * 7
)
