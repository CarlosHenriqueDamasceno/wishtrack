package register

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type RegisterHandler struct {
}

type RegisterOutput struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	HandleRegister(w, req)
}

func HandleRegister(w http.ResponseWriter, req *http.Request) {
	out := &RegisterOutput{}
	err := json.NewDecoder(req.Body).Decode(out)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	out.ID = uuid.New()
	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
