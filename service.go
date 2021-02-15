package main

import (
	"fmt"
	"time"
)

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

	cachedTranslator, err := newCachedTranslator(&backoffTranslator{
		translator: t,
		retries:    3,
		backoff:    time.Millisecond,
	}, 100)
	if err != nil {
		panic(fmt.Errorf("newCachedTranslator failed with error: %v", err))
	}

	return &Service{
		translator: cachedTranslator,
	}
}
