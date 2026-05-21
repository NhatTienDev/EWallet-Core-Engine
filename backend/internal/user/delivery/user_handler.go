package delivery

import (
	"github.com/nhattiendev/ewallet/internal/user/domain"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userUC domain.UserUseCase
}

func NewUserHandler(userUC domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUC: userUC,
	}
}	

// Declare API routes for User module
func (h *UserHandler) RegisterUserRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("api/v1/users", func(r chi.Router) {
		r.Post("/register", h.HandleRegister)
		r.Post("/login", h.HandleLogin)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/profile", h.HandleGetProfile)
		})
	})
}