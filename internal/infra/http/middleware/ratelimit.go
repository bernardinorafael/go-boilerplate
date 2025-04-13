package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"golang.org/x/time/rate"
)

const (
	// rateLimit is the number of requests per second
	rateLimit = 2
	// burst is the maximum number of requests that can be made in a single burst
	burst = 4
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func withRateLimit(next http.Handler) http.Handler {
	var mu sync.Mutex
	var clients = make(map[string]*client)

	// Background routine to remove expired clients
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Minute*3 {
					delete(clients, ip)
				}
			}
			// IMPORTANT: Unlock the mutext when the cleanup is done
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()

		ip := r.RemoteAddr
		if _, ok := clients[ip]; !ok {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rateLimit, burst),
			}
		}
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			fault.NewHTTPError(w, fault.New(
				"too many requests",
				fault.WithHTTPCode(http.StatusTooManyRequests),
				fault.WithTag(fault.TOO_MANY_REQUESTS),
			))
			return
		}

		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
