package sinks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/resmoio/kubernetes-event-exporter/pkg/kube"
	"github.com/stretchr/testify/assert"
)

func TestLoki_SendWithBasicAuth(t *testing.T) {
	username := "testuser"
	password := "testpass"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		expectedAuth := "Basic " + basicAuth(username, password)
		assert.Equal(t, expectedAuth, auth)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := &LokiConfig{
		URL:      ts.URL,
		Username: username,
		Password: password,
	}

	loki, err := NewLoki(cfg)
	assert.NoError(t, err)

	ev := &kube.EnhancedEvent{}
	err = loki.Send(context.Background(), ev)
	assert.NoError(t, err)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestLoki_SendWithoutBasicAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		assert.Empty(t, auth)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := &LokiConfig{
		URL: ts.URL,
	}

	loki, err := NewLoki(cfg)
	assert.NoError(t, err)

	ev := &kube.EnhancedEvent{}
	err = loki.Send(context.Background(), ev)
	assert.NoError(t, err)
}
