package http

import (
	"net/http"

	"ferreteria/internal/config"
)

// Category endpoints exposed by this router:
//   GET /api/categories   list categories, optionally filtered by type
//   POST /api/categories  create a custom category
func registerCategoryRoutes(
	mux *http.ServeMux,
	settings config.Config,
	dependencies *dependencies,
) {
	collection := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			handleListCategory(dependencies)(writer, request)
		case http.MethodPost:
			handleCreateCategory(dependencies)(writer, request)
		default:
			writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	_ = settings
	mux.Handle("/api/categories", withAuth(dependencies.tokens, collection))
}
