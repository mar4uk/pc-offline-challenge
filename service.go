package main

import "time"

// Service is a Translator user.
type Service struct {
	translator Translator
}

// NewService is a constuctor for translation service
func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)

	return &Service{
		translator: &backoffTranslator{
			translator: t,
			retries:    3,
			backoff:    time.Millisecond,
		},
	}
}
