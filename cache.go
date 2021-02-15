package main

import (
	"context"
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/text/language"
)

type cachedTranslator struct {
	translator Translator
	cache      *lru.Cache
}

func getKeyForCache(from, to language.Tag, data string) string {
	return fmt.Sprintf("%s-%s-%s", from, to, data)
}

func newCachedTranslator(translator Translator, size int) (Translator, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	return &cachedTranslator{
		translator: translator,
		cache:      cache,
	}, nil
}

func (t *cachedTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (result string, err error) {
	key := getKeyForCache(from, to, data)

	value, ok := t.cache.Get(key)
	if ok {
		return value.(string), nil
	}

	result, err = t.translator.Translate(ctx, from, to, data)
	if err != nil {
		return "", err
	}

	t.cache.Add(key, result)
	return result, nil
}
