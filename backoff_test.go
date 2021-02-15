package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
)

type mockedTranslator struct {
	mock.Mock
}

func (t *mockedTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	args := t.Called(ctx, from, to, data)
	return args.String(0), args.Error(1)
}

func TestBackoffNoError(t *testing.T) {
	mocked := new(mockedTranslator)
	translator := &backoffTranslator{
		translator: mocked,
	}

	ctx := context.TODO()
	from := language.Afrikaans
	to := language.English
	testIn := "some text"
	testOut := "some result"

	mocked.On("Translate", ctx, from, to, testIn).Return(testOut, nil).Once()

	res, err := translator.Translate(ctx, from, to, testIn)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, testOut)

	mocked.AssertExpectations(t)
}

func TestBackoffError(t *testing.T) {
	mocked := new(mockedTranslator)
	translator := &backoffTranslator{
		translator: mocked,
		retries:    1,
	}

	ctx := context.TODO()
	from := language.Afrikaans
	to := language.English
	testIn := "some text"
	testError := fmt.Errorf("test error")

	mocked.On("Translate", ctx, from, to, testIn).Return("", testError).Twice()

	res, err := translator.Translate(ctx, from, to, testIn)
	assert.Equal(t, err, testError)
	assert.Empty(t, res)

	mocked.AssertExpectations(t)
}

func TestBackoff1Retry(t *testing.T) {
	mocked := new(mockedTranslator)
	translator := &backoffTranslator{
		translator: mocked,
		retries:    1,
		backoff:    100 * time.Millisecond,
	}

	ctx := context.TODO()
	from := language.Afrikaans
	to := language.English
	testIn := "some text"
	testOut := "some result"
	testError := fmt.Errorf("test error")

	var firstCallTime time.Time
	mocked.On("Translate", ctx, from, to, testIn).Return("", testError).Once().Run(func(_ mock.Arguments) {
		firstCallTime = time.Now()
	})
	var secondCallTime time.Time
	mocked.On("Translate", ctx, from, to, testIn).Return(testOut, nil).Once().Run(func(_ mock.Arguments) {
		secondCallTime = time.Now()
	})

	start := time.Now()
	res, err := translator.Translate(ctx, from, to, testIn)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, testOut)

	assert.InDelta(t, firstCallTime.Sub(start), 0*time.Millisecond, float64(time.Millisecond))
	assert.InEpsilon(t, secondCallTime.Sub(firstCallTime), 100*time.Millisecond, 0.1)
	mocked.AssertExpectations(t)
}
