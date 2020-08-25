package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestClient_MakeRequest(t *testing.T) {
	defer gock.Off() // flush pending mocks after test execution

	gock.New("https://snoo-api.happiestbaby.com").
		Get("/some/path").
		Reply(200).
		JSON(map[string]string{"status": "ok"})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewClient("username", "password")

	// declare an empty interface
	var result map[string]interface{}

	client.MakeRequest(ctx, http.MethodGet, "/some/path", nil, &result)

	assert.Equal(t, "ok", result["status"])
}
