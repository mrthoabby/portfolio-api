package main

import "time"

const (
	// Security constants
	maxBodySize        = 1 << 20 // 1 MB
	readTimeout        = 10 * time.Second
	writeTimeout       = 30 * time.Second
	idleTimeout        = 120 * time.Second
	shutdownTimeout    = 10 * time.Second
	rateLimitRequests  = 100             // requests per window
	rateLimitWindow    = 1 * time.Minute // time window
	contactRateLimit   = 5               // contact requests per window
	contactRateWindow  = 1 * time.Minute // time window for contacts
	questionRateLimit  = 10              // question requests per window
	questionRateWindow = 1 * time.Minute // time window for questions
)
