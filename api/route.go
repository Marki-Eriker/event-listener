package api

import (
	"github.com/marki-eriker/event-listener/entity/user"
	"net/http"
)

var AdminOnly = map[user.Role]struct{}{user.Admin: {}}
var AnalystOnly = map[user.Role]struct{}{user.Analyst: {}}

func (s *Server) initMainRoutesV1() {
	v1 := s.router.PathPrefix("/v1").Subrouter()

	authRouter := v1.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", s.handleLogin).Methods(http.MethodPost)
	authRouter.HandleFunc("/register", s.handleRegisterUser).Methods(http.MethodPost)

	userRouter := v1.PathPrefix("/user").Subrouter()
	userRouter.Use(authMiddleware(s.tokenManager))

	verifyRouter := userRouter.PathPrefix("/verify").Subrouter()
	verifyRouter.Use(withRoleMiddleware(AdminOnly))
	verifyRouter.HandleFunc("/{id}", s.handleVerifyUser).Methods(http.MethodPatch)

	eventRouter := v1.PathPrefix("/event").Subrouter()
	eventRouter.Use(authMiddleware(s.tokenManager))

	eventRouter.HandleFunc("/list", s.handleGetEvents).Methods(http.MethodGet)

	incidentRouter := eventRouter.PathPrefix("/incident").Subrouter()
	incidentRouter.Use(withRoleMiddleware(AnalystOnly))
	incidentRouter.HandleFunc("/{id}", s.handleMarkEventAnIncident).Methods(http.MethodPatch)
}
