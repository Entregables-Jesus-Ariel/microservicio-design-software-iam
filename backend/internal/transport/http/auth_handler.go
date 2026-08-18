package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"iam/internal/application/usecase"
)

// handleRegister creates a new user account with the default role.
func handleRegister(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload registerRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		if strings.TrimSpace(payload.Email) == "" || len(payload.Password) < 8 {
			writeError(w, http.StatusBadRequest, "email is required and password must have at least 8 characters")
			return
		}

		user, err := deps.registerUser.Execute(r.Context(), registerInput(payload))
		if err != nil {
			writeDomainError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, userResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}
}

func registerInput(p registerRequest) usecase.RegisterUserInput {
	return usecase.RegisterUserInput{
		Email:     p.Email,
		Password:  p.Password,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}
}

// handleLogin authenticates a user and issues a token pair on success.
func handleLogin(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload loginRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		if strings.TrimSpace(payload.Email) == "" || payload.Password == "" {
			writeError(w, http.StatusBadRequest, "email and password are required")
			return
		}

		result, err := deps.loginUser.Execute(r.Context(), usecase.LoginUserInput{
			Email:    payload.Email,
			Password: payload.Password,
		})
		if err != nil {
			writeDomainError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, loginResponse{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			User: userResponse{
				ID:        result.User.ID,
				Email:     result.User.Email,
				FirstName: result.User.FirstName,
				LastName:  result.User.LastName,
			},
		})
	}
}

// handleRefresh issues a new token pair from a valid refresh token.
func handleRefresh(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload refreshRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		if strings.TrimSpace(payload.RefreshToken) == "" {
			writeError(w, http.StatusBadRequest, "refresh token is required")
			return
		}

		result, err := deps.refreshSession.Execute(r.Context(), usecase.RefreshSessionInput{
			RefreshToken: payload.RefreshToken,
		})
		if err != nil {
			writeDomainError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, refreshResponse{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
		})
	}
}

// handleLogout invalidates a refresh token so it can no longer be used.
func handleLogout(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload logoutRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		if strings.TrimSpace(payload.RefreshToken) != "" {
			_ = deps.logoutUser.Execute(r.Context(), usecase.LogoutUserInput{
				RefreshToken: payload.RefreshToken,
			})
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// handleForgotPassword requests a password reset token.
func handleForgotPassword(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload forgotPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		if strings.TrimSpace(payload.Email) == "" {
			writeError(w, http.StatusBadRequest, "email is required")
			return
		}

		result, err := deps.forgotPassword.Execute(r.Context(), usecase.ForgotPasswordInput{
			Email: payload.Email,
		})
		if err != nil {
			writeDomainError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, forgotPasswordResponse{
			ResetToken: result.ResetToken,
		})
	}
}

// handleResetPassword sets a new password using a valid reset token.
func handleResetPassword(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload resetPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		if strings.TrimSpace(payload.Token) == "" || len(payload.NewPassword) < 8 {
			writeError(w, http.StatusBadRequest, "valid token and new password (min 8 chars) are required")
			return
		}

		if err := deps.resetPassword.Execute(r.Context(), usecase.ResetPasswordInput{
			Token:       payload.Token,
			NewPassword: payload.NewPassword,
		}); err != nil {
			writeDomainError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
