package main

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/text/language"
)

type deduplicatedTranslator struct {
	translator Translator
	set        map[string]struct{}
	cond       *sync.Cond
}

func newDeduplicatedTranslator(translator Translator) Translator {
	return &deduplicatedTranslator{
		translator: translator,
		set:        make(map[string]struct{}),
		cond:       sync.NewCond(&sync.Mutex{}),
	}
}

func (t *deduplicatedTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (result string, err error) {
	key := fmt.Sprintf("%s-%s-%s", from, to, data)

	condition := func() bool {
		_, ok := t.set[key]
		return ok
	}

	t.cond.L.Lock()
	for condition() {
		t.cond.Wait()
	}

	t.set[key] = struct{}{}
	t.cond.L.Unlock()

	result, err = t.translator.Translate(ctx, from, to, data)
	if err != nil {
		return "", err
	}

	t.cond.L.Lock()
	delete(t.set, key)
	t.cond.Broadcast()
	t.cond.L.Unlock()

	return result, err
}
