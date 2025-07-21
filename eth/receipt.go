package eth

import (
	"time"

	"golang.org/x/time/rate"
)

var ReceiptRateLimiter *rate.Limiter

func InitReceiptRateLimiter() {
	ReceiptRateLimiter = rate.NewLimiter(rate.Every(200*time.Millisecond), 1)
}
