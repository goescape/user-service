package config

import (
	"time"

	"github.com/sony/gobreaker"
)

func InitBreaker() *gobreaker.Settings {
	return &gobreaker.Settings{
		Name:        "ProductServiceBreaker",
		MaxRequests: 3,
		Interval:    time.Second * 60,
		Timeout:     time.Second * 10,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	}
}
