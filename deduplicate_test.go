package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
)

type mockedDeduplicateTranslator struct {
	mock.Mock
}

func (t *mockedDeduplicateTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	args := t.Called(ctx, from, to, data)
	time.Sleep(100 * time.Millisecond)
	return args.String(0), args.Error(1)
}

func TestDeduplicate(t *testing.T) {
	mocked := new(mockedDeduplicateTranslator)
	translator := newDeduplicatedTranslator(mocked)

	ctx := context.TODO()
	from := language.Afrikaans
	to := language.English
	testIn := "some text"
	testOut := "some result"
	wg := &sync.WaitGroup{}

	var firstCallTime time.Time
	mocked.On("Translate", ctx, from, to, testIn).Return(testOut, nil).Once().Run(func(_ mock.Arguments) {
		firstCallTime = time.Now()
	})
	var secondCallTime time.Time
	mocked.On("Translate", ctx, from, to, testIn).Return(testOut, nil).Once().Run(func(_ mock.Arguments) {
		secondCallTime = time.Now()
	})

	for i := 0; i < 2; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			translator.Translate(ctx, from, to, testIn)
			wg.Done()
		}(wg)
	}

	wg.Wait()

	assert.InEpsilon(t, secondCallTime.Sub(firstCallTime), 100*time.Millisecond, 0.1)

	mocked.AssertExpectations(t)
}
