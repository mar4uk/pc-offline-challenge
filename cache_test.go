package main

import (
	"context"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestHasCacheForQuery(t *testing.T) {
	mocked := new(mockedTranslator)

	ctx := context.TODO()
	from := language.Afrikaans
	to := language.English
	testIn := "some text"
	testOut := "some result"

	mockedCache, _ := lru.New(10)
	mockedCache.Add(getKeyForCache(from, to, testIn), testOut)

	translator := &cachedTranslator{
		translator: mocked,
		cache:      mockedCache,
	}

	res, err := translator.Translate(ctx, from, to, testIn)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, testOut)

	mocked.AssertNotCalled(t, "Translate", ctx, from, to, testIn)
}

func TestHasNoCacheForQuery(t *testing.T) {
	mocked := new(mockedTranslator)
	mockedCache, _ := lru.New(10)

	translator := &cachedTranslator{
		translator: mocked,
		cache:      mockedCache,
	}

	ctx := context.TODO()
	from := language.Afrikaans
	to := language.English
	testIn := "some text"
	testOut := "some result"

	mocked.On("Translate", ctx, from, to, testIn).Return(testOut, nil).Once()

	res, err := translator.Translate(ctx, from, to, testIn)
	value, ok := mockedCache.Get(getKeyForCache(from, to, testIn))

	assert.Equal(t, err, nil)
	assert.Equal(t, res, testOut)
	assert.Equal(t, ok, true)
	assert.Equal(t, value.(string), testOut)

	mocked.AssertExpectations(t)
}
