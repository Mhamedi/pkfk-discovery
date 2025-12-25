package middleware

import (
	"net/http"

	"github.com/pkfk-discovery/api/internal/domain"
)

type RBACMiddleware struct {
	requiredRole domain.Role
}

func NewRBACMiddleware(requiredRole domain.Role) *RBACMiddleware {
	return &RBACMiddleware{
		requiredRole: requiredRole,
	}
}

func (m *RBACMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetUserRole(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userRole := domain.Role(role)
		if !hasPermission(userRole, m.requiredRole) {
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func hasPermission(userRole, requiredRole domain.Role) bool {
	if userRole == domain.RoleAdmin {
		return true
	}

	if requiredRole == domain.RoleViewer {
		return userRole == domain.RoleViewer || userRole == domain.RoleEditor
	}

	if requiredRole == domain.RoleEditor {
		return userRole == domain.RoleEditor
	}

	return false
}

