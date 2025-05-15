package pluggy

import (
	"sync"
	"time"
)

const maxRequests = 360
const rateLimitDuration = time.Hour

type rateLimiter struct {
	mu        sync.Mutex
	timestamp time.Time
	count     int
}

func (rl *rateLimiter) wait() {
	for {
		rl.mu.Lock()
		now := time.Now()

		// Check if the rate limit duration has passed
		if now.Sub(rl.timestamp) > rateLimitDuration {
			// Reset the timestamp and request count
			rl.timestamp = now
			rl.count = 0
		}

		// Allow the request if the count is below the maximum allowed requests
		if rl.count < maxRequests {
			rl.count++
			rl.mu.Unlock()
			return
		}

		rl.mu.Unlock()

		// Calculate the sleep time before retrying
		sleepTime := rateLimitDuration - now.Sub(rl.timestamp)
		if sleepTime > time.Second {
			sleepTime = time.Second
		}

		time.Sleep(sleepTime)
	}
}
