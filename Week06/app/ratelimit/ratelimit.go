package ratelimit


type RateLimit interface {
	Take() error
}
