package http

import (
	"net/http"
)

type AuthenticatedTransport struct {
	T     http.RoundTripper
	token string
}

func NewAuthenticatedTransport(token string) *AuthenticatedTransport {
	return &AuthenticatedTransport{
		T:     http.DefaultTransport.(*http.Transport).Clone(),
		token: "Bearer " + token,
	}
}

func (t *AuthenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", t.token)
	return t.T.RoundTrip(req)
}
