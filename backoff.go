package main

import (
	"context"
	"math"
	"time"

	"golang.org/x/text/language"
)

type backoffTranslator struct {
	translator Translator
	retries    uint
	backoff    time.Duration
}

func (t *backoffTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (result string, err error) {
	for retry := uint(0); retry <= t.retries; retry++ {
		result, err = t.translator.Translate(ctx, from, to, data)

		if err == nil {
			return result, err
		}

		// increase sleep time exponentially for the next retry
		time.Sleep(time.Duration(math.Pow(float64(2), float64(retry)) * float64(t.backoff)))
	}

	return result, err
}
