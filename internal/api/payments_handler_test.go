package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"komodo-customer-api/internal/db"
	"komodo-customer-api/internal/models"
)

// ── Unit Tests: GetPaymentsHandler ───────────────────────────────────────────

func TestGetPaymentsHandler_TokenRedacted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	repo.EXPECT().
		ListPayments(gomock.Any(), "user_abc").
		Return([]models.PaymentMethod{
			{PaymentID: "pay_1", Token: "secret-token", Last4: "4242"},
		}, nil)

	req := makeRequest(t, http.MethodGet, "/v1/me/payments", nil)
	req = withUserID(req, "user_abc")
	rr := httptest.NewRecorder()
	svc.GetPaymentsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotContains(t, rr.Body.String(), "secret-token")
	assert.Contains(t, rr.Body.String(), "pay_1")
}

func TestGetPaymentsHandler_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	req := makeRequest(t, http.MethodGet, "/v1/me/payments", nil)
	rr := httptest.NewRecorder()
	svc.GetPaymentsHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// ── Unit Tests: UpsertPaymentHandler ─────────────────────────────────────────

func TestUpsertPaymentHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	repo.EXPECT().
		UpsertPayment(gomock.Any(), "user_abc", gomock.Any()).
		Return(nil)

	body := map[string]any{"payment_id": "pay_1", "provider": "stripe", "last4": "4242"}
	req := withUserID(makeRequest(t, http.MethodPut, "/v1/me/payments", body), "user_abc")
	rr := httptest.NewRecorder()
	svc.UpsertPaymentHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpsertPaymentHandler_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	req, _ := http.NewRequest(http.MethodPut, "/v1/me/payments", strings.NewReader("not-json"))
	req = withUserID(req, "user_abc")
	rr := httptest.NewRecorder()
	svc.UpsertPaymentHandler(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpsertPaymentHandler_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, _ := newTestService(t, ctrl)

	body := map[string]any{"payment_id": "pay_1"}
	req := makeRequest(t, http.MethodPut, "/v1/me/payments", body)
	rr := httptest.NewRecorder()
	svc.UpsertPaymentHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// ── Unit Tests: DeletePaymentHandler ─────────────────────────────────────────

func TestDeletePaymentHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	repo.EXPECT().
		DeletePayment(gomock.Any(), "user_abc", "pay_1").
		Return(nil)

	req := withUserID(makeRequest(t, http.MethodDelete, "/v1/me/payments/pay_1", nil), "user_abc")
	req.SetPathValue("id", "pay_1")
	rr := httptest.NewRecorder()
	svc.DeletePaymentHandler(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeletePaymentHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc, repo := newTestService(t, ctrl)

	repo.EXPECT().
		DeletePayment(gomock.Any(), "user_abc", "pay_missing").
		Return(db.ErrNotFound)

	req := withUserID(makeRequest(t, http.MethodDelete, "/v1/me/payments/pay_missing", nil), "user_abc")
	req.SetPathValue("id", "pay_missing")
	rr := httptest.NewRecorder()
	svc.DeletePaymentHandler(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
