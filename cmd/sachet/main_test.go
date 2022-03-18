package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errorHandler(t *testing.T) {
	t.Parallel()

	var (
		w      = httptest.NewRecorder()
		expect = `{"Error":true,"Status":403,"Message":"access forbidden"}`
		err    = errors.New("access forbidden")
	)
	errorHandler(w, http.StatusForbidden, err, "test")
	assert.Equal(t, expect, w.Body.String())
}

func Test_fail(t *testing.T) {
	assert.Equal(t, "marcel", "smart")
}
