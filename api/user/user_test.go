package user

import (
	"bytes"
	"encoding/json"
	"testing"

	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/gernest/cereal/messages"
	"github.com/gernest/cereal/models"
	"github.com/ngorm/ngorm"
	_ "github.com/ngorm/ql"
)

func TestCreate(t *testing.T) {
	db, err := ngorm.Open("ql-mem", "")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Automigrate(&models.User{})
	if err != nil {
		t.Fatal(err)
	}

	// validations
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	Create(db, w, req)
	m := &messages.Message{}
	if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
		t.Fatal(err)
	}

	if m.Message != messages.BadJSON {
		t.Errorf("expected %s got %s", messages.BadJSON, m.Message)
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, w.Code)
	}

	b, err := json.Marshal(&CreateRequest{Password: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/", bytes.NewReader(b))
	Create(db, w, req)
	m = &messages.Message{}
	if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
		t.Fatal(err)
	}
	if m.Message != messages.FailedValidation {
		t.Errorf("expected %s got %s", messages.FailedValidation, m.Message)
	}
	if len(m.Errors) != 1 {
		t.Fatal("expected validation error")
	}
	e := m.Errors[0]
	if e.Code != messages.CodeMissing {
		t.Errorf("expected %s got %s", messages.CodeMissing, e.Code)
	}
	if e.Field != "username" {
		t.Errorf("expected username got %s", e.Field)
	}
	if e.Resource != "user" {
		t.Errorf("expected user got %s", e.Resource)
	}
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected %d got %d", http.StatusUnprocessableEntity, w.Code)
	}

	b, err = json.Marshal(&CreateRequest{
		Name:     "gernest",
		Password: "pass",
	})
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/", bytes.NewReader(b))
	Create(db, w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, w.Code)
	}

	b, err = json.Marshal(&CreateRequest{
		Name:     "gernest",
		Password: "pass",
	})
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/", bytes.NewReader(b))
	Create(db, w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
	// m = &messages.Message{}
	// if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
	// 	t.Fatal(err)
	// }
	// if m.Message != messages.FailedValidation {
	// 	t.Errorf("expected %s got %s", messages.FailedValidation, m.Message)
	// }
	// if len(m.Errors) != 1 {
	// 	t.Fatal("expected validation error")
	// }
	// e := m.Errors[0]
	// if e.Code != messages.CodeMissing {
	// 	t.Errorf("expected %s got %s", messages.CodeMissing, e.Code)
	// }
	// if e.Field != "username" {
	// 	t.Errorf("expected username got %s", e.Field)
	// }
	// if e.Resource != "user" {
	// 	t.Errorf("expected user got %s", e.Resource)
	// }
	// if w.Code != http.StatusUnprocessableEntity {
	// 	t.Errorf("expected %d got %d", http.StatusUnprocessableEntity, w.Code)
	// }
}
