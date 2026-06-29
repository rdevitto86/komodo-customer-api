package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// ── Unit Tests: UnsubscribeHandler ───────────────────────────────────────────

func TestUnsubscribeHandler_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	token, err := svc.MintUnsubscribeToken(context.Background(), "cust_1", "email")
	require.NoError(t, err)

	repo.EXPECT().AppendConsentLog(gomock.Any(), "cust_1", gomock.Any()).Return(nil)

	req := makeRequest(t, http.MethodPost, "/v1/unsubscribe", map[string]any{"token": token})
	rr := httptest.NewRecorder()
	svc.UnsubscribeHandler(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestUnsubscribeHandler_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	req := makeRequest(t, http.MethodPost, "/v1/unsubscribe", map[string]any{"token": "invalid-token"})
	rr := httptest.NewRecorder()
	svc.UnsubscribeHandler(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// ── Unit Tests: MintUnsubscribeTokenHandler ───────────────────────────────────

func TestMintUnsubscribeTokenHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	body := map[string]any{"channel": "email"}
	req := makeRequest(t, http.MethodPost, "/v1/customers/cust_1/unsubscribe-token", body)
	req.SetPathValue("id", "cust_1")
	rr := httptest.NewRecorder()
	svc.MintUnsubscribeTokenHandler(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "token")
}

func TestMintUnsubscribeTokenHandler_NoPathID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	body := map[string]any{"channel": "email"}
	req := makeRequest(t, http.MethodPost, "/v1/customers//unsubscribe-token", body)
	rr := httptest.NewRecorder()
	svc.MintUnsubscribeTokenHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
