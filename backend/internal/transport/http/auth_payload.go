package http

// loginRequest carries the credentials submitted by the admin.
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// tokenResponse returns the issued access token.
type tokenResponse struct {
	Token string `json:"token"`
}
