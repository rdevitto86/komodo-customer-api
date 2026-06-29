package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ── Unit Tests: ExportProfileHandler ─────────────────────────────────────────

func TestExportProfileHandler_NoS3Client(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	req := withUserID(makeRequest(t, http.MethodPost, "/v1/me/profile/export", nil), "user_abc")
	rr := httptest.NewRecorder()
	svc.ExportProfileHandler(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestExportProfileHandler_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	req := makeRequest(t, http.MethodPost, "/v1/me/profile/export", nil)
	rr := httptest.NewRecorder()
	svc.ExportProfileHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
