package http

import (
	"net/http"
	"strings"

	"ferreteria/internal/config"
)

// Movement endpoints exposed by this router:
//   POST /api/movements              record an income or expense entry
//   GET /api/movements               list movements of a period
//   PUT /api/movements/{id}          edit an existing movement
//   DELETE /api/movements/{id}       cancel a movement as a soft delete
//   GET /api/movements/{id}/audit    list the change history of a movement
func registerMovementRoutes(
	mux *http.ServeMux,
	settings config.Config,
	dependencies *dependencies,
) {
	collection := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			handleCreateMovement(dependencies)(writer, request)
		case http.MethodGet:
			handleListMovement(dependencies)(writer, request)
		default:
			writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	item := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path := strings.TrimSuffix(request.URL.Path, "/")
		if strings.HasSuffix(path, "/audit") {
			if request.Method != http.MethodGet {
				writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleGetMovementAudit(dependencies)(writer, request)
			return
		}

		switch request.Method {
		case http.MethodPut:
			handleEditMovement(dependencies)(writer, request)
		case http.MethodDelete:
			handleCancelMovement(dependencies)(writer, request)
		default:
			writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	_ = settings
	mux.Handle("/api/movements", withAuth(dependencies.tokens, collection))
	mux.Handle("/api/movements/", withAuth(dependencies.tokens, item))
}