package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type T struct {
	*testing.T
}

func TestUnkownRouteReturns404(t *testing.T) {
	recorder := requestAPI(t, "GET", "/", nil)
	assert.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
}

func TestHydrateEndpointReturns200(t *testing.T) {
	recorder := requestAPI(t, "GET", "/hydrate", nil)
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
}

func TestHydrateEndpointReturnsMessage(t *testing.T) {
	recorder := requestAPI(t, "GET", "/hydrate", nil)
	assert.Equal(t, "Let's hydrate!", recorder.Body.String())
}

func requestAPI(t *testing.T, method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := BuildRouter()
	router.ServeHTTP(rr, req)
	return rr
}
