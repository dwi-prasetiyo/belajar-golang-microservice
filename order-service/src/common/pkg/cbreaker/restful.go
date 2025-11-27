package cbreaker

import (
	"fmt"
	"net/http"
	"order-service/src/common/errors"
	"order-service/src/common/log"
	"slices"
	"time"

	"github.com/sony/gobreaker/v2"
)

func NewRestful(name string) *gobreaker.CircuitBreaker[any] {
	setting := gobreaker.Settings{
		Name:        name,
		MaxRequests: 5,
		Interval:    1 * time.Minute,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRate := float32(counts.TotalFailures) / float32(counts.Requests)
			return failureRate >= 0.8 && counts.Requests >= 5
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}

			if e, ok := err.(*errors.Response); ok {
				statusCodeFailed := []int{
					http.StatusInternalServerError,
					http.StatusBadGateway,
					http.StatusServiceUnavailable,
					http.StatusGatewayTimeout,
					599,
				}
				return !slices.Contains(statusCodeFailed, e.HttpCode)
			}

			return false
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Logger.Info(fmt.Sprintf("[%s] state change: %s -> %s", name, from, to))
		},
	}

	return gobreaker.NewCircuitBreaker[any](setting)
}
