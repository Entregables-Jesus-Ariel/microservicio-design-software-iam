package http

import (
	"encoding/json"
	"net/http"

	"iam/internal/application/port"
	"iam/internal/application/usecase"
)

func handleListUsers(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usersWithRoles, err := deps.rbacUsecases.ListUsersWithRoles(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list users")
			return
		}
		writeJSON(w, http.StatusOK, usersWithRoles)
	}
}

func handleListRoles(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles, err := deps.rbacUsecases.ListRoles(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list roles")
			return
		}
		writeJSON(w, http.StatusOK, roles)
	}
}

type assignRoleRequest struct {
	RoleID string `json:"roleId"`
}

func handleAssignRole(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")
		if userID == "" {
			writeError(w, http.StatusBadRequest, "user id is required")
			return
		}

		var payload assignRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		claims, _ := r.Context().Value(userClaimsKey).(port.TokenClaims)
		assignedBy := claims.UserID // Ensure we use the token's user ID

		err := deps.rbacUsecases.AssignRole(r.Context(), usecase.AssignRoleInput{
			UserID:     userID,
			RoleID:     payload.RoleID,
			AssignedBy: assignedBy,
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to assign role")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handleRevokeRole(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")
		roleID := r.PathValue("roleId")

		if userID == "" || roleID == "" {
			writeError(w, http.StatusBadRequest, "user id and role id are required")
			return
		}

		err := deps.rbacUsecases.RevokeRole(r.Context(), userID, roleID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to revoke role")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
