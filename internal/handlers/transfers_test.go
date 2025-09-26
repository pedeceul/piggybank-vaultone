package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTransfer_Validation(t *testing.T) {
	body := bytes.NewBufferString(`{"from_account_id":"a1","amount":"10","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/transfers", body)
	rr := httptest.NewRecorder()
	CreateTransfer(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestCreateTransfer_Missing(t *testing.T) {
	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/transfers", body)
	rr := httptest.NewRecorder()
	CreateTransfer(rr, req)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
}
