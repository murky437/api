package clock

import "time"

type Clock interface {
	Now() time.Time
}

type defaultClock struct{}

func (c defaultClock) Now() time.Time {
	return time.Now()
}

func New() Clock {
	return defaultClock{}
}
