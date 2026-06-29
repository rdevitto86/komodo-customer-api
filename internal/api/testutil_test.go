package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	ctxKeys "github.com/rdevitto86/komodo-forge-sdk-go/http/context"
	"go.uber.org/mock/gomock"

	"komodo-customer-api/internal/api/mocks"
)

func newTestService(t *testing.T, ctrl *gomock.Controller) (*Service, *mocks.Mockrepository) {
	t.Helper()
	repo := mocks.NewMockrepository(ctrl)
	svc := NewService(repo, ServiceExtraConfig{UnsubscribeKey: []byte("test-secret-32-bytes-padded-xx!!")})
	return svc, repo
}

func makeRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()
	if body == nil {
		req, err := http.NewRequest(method, path, http.NoBody)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		return req
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}
	req, err := http.NewRequest(method, path, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func withUserID(req *http.Request, id string) *http.Request {
	ctx := context.WithValue(req.Context(), ctxKeys.USER_ID_KEY, id)
	return req.WithContext(ctx)
}

func withScopes(req *http.Request, scopes []string) *http.Request {
	ctx := context.WithValue(req.Context(), ctxKeys.SCOPES_KEY, scopes)
	return req.WithContext(ctx)
}
