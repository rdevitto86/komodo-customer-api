package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"komodo-customer-api/internal/db"
	"komodo-customer-api/internal/models"
)

// ── Unit Tests: GetSettingsHandler ───────────────────────────────────────────

func TestGetSettingsHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	repo.EXPECT().
		GetSettings(gomock.Any(), "cust_1").
		Return(&models.AccountSettings{Status: "active"}, nil)

	req := makeRequest(t, http.MethodGet, "/v1/customers/cust_1/settings", nil)
	req.SetPathValue("id", "cust_1")
	rr := httptest.NewRecorder()
	svc.GetSettingsHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetSettingsHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	repo.EXPECT().
		GetSettings(gomock.Any(), "cust_missing").
		Return(nil, db.ErrNotFound)

	req := makeRequest(t, http.MethodGet, "/v1/customers/cust_missing/settings", nil)
	req.SetPathValue("id", "cust_missing")
	rr := httptest.NewRecorder()
	svc.GetSettingsHandler(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetSettingsHandler_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	req := makeRequest(t, http.MethodGet, "/v1/customers//settings", nil)
	rr := httptest.NewRecorder()
	svc.GetSettingsHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// ── Unit Tests: UpdateSettingsHandler ────────────────────────────────────────

func TestUpdateSettingsHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	repo.EXPECT().
		UpdateSettings(gomock.Any(), "cust_1", gomock.Any()).
		Return(nil)

	body := map[string]any{"status": "active"}
	req := makeRequest(t, http.MethodPut, "/v1/customers/cust_1/settings", body)
	req.SetPathValue("id", "cust_1")
	rr := httptest.NewRecorder()
	svc.UpdateSettingsHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateSettingsHandler_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	req := makeRequest(t, http.MethodPut, "/v1/customers/cust_1/settings", nil)
	req.SetPathValue("id", "cust_1")
	rr := httptest.NewRecorder()
	svc.UpdateSettingsHandler(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateSettingsHandler_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	req := makeRequest(t, http.MethodPut, "/v1/customers//settings", map[string]any{"status": "active"})
	rr := httptest.NewRecorder()
	svc.UpdateSettingsHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// ── Unit Tests: UpdateSettingsTagsHandler ────────────────────────────────────

func TestUpdateSettingsTagsHandler_ValidNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	repo.EXPECT().
		GetSettings(gomock.Any(), "cust_1").
		Return(&models.AccountSettings{Status: "active"}, nil)
	repo.EXPECT().
		UpdateSettings(gomock.Any(), "cust_1", gomock.Any()).
		Return(nil)

	body := map[string]any{"add": []string{"loyalty.vip"}, "remove": []string{}}
	req := makeRequest(t, http.MethodPut, "/v1/customers/cust_1/settings/tags", body)
	req.SetPathValue("id", "cust_1")
	req = withScopes(req, []string{"svc:loyalty-api"})
	rr := httptest.NewRecorder()
	svc.UpdateSettingsTagsHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateSettingsTagsHandler_ForbiddenNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	body := map[string]any{"add": []string{"system.internal"}, "remove": []string{}}
	req := makeRequest(t, http.MethodPut, "/v1/customers/cust_1/settings/tags", body)
	req.SetPathValue("id", "cust_1")
	req = withScopes(req, []string{"svc:loyalty-api"})
	rr := httptest.NewRecorder()
	svc.UpdateSettingsTagsHandler(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}
