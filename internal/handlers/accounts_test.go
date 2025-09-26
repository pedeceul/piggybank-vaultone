package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAccount_Validation(t *testing.T) {
	body := bytes.NewBufferString(`{"owner_id":"u1","kind":"checking","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/accounts", body)
	rr := httptest.NewRecorder()
	CreateAccount(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestCreateAccount_Missing(t *testing.T) {
	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/accounts", body)
	rr := httptest.NewRecorder()
	CreateAccount(rr, req)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
}
