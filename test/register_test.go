package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CarlosHenriqueDamasceno/wishtrack/api"
	"github.com/google/uuid"
)

func TestShouldRegister(t *testing.T) {
	input := struct {
		Name     string
		Email    string
		Password string
	}{
		Name:     "Carlos",
		Email:    "carlos@wishtrack.com",
		Password: "password",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Error("Body should be serialized")
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))

	handler := api.NewApiServer(http.NewServeMux())
	handler.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != http.StatusCreated {
		t.Errorf("Status code should be 201 created, received: %s", recorder.Result().Status)
	}

	responseBody := struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}

	err = json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	if err != nil {
		t.Errorf("Fail to unmarshal response: %s", err.Error())
	}
	err = uuid.Validate(responseBody.ID)
	if err != nil {
		t.Errorf("Result ID is invalid: %s", err.Error())
	}
}
