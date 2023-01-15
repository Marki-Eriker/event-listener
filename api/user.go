package api

import (
	"encoding/json"
	"github.com/marki-eriker/event-listener/entity/user"
	"github.com/marki-eriker/event-listener/pkg/request"
	"github.com/marki-eriker/event-listener/processor"
	"net/http"
)

type registerUserInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     uint16 `json:"role"`
}

func (s *Server) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	req := request.New(w, r)

	if req.GetBody() == nil || len(req.GetBody()) == 0 {
		req.FinishBadRequest("No JSON body")
		return
	}

	var payload registerUserInput

	err := json.Unmarshal(req.GetBody(), &payload)
	if err != nil {
		req.FinishBadRequest("Unable to unmarshal JSON data: %s", err)
		return
	}

	err = s.processor.RegisterUser(r.Context(), payload.Login, payload.Password, payload.Role)
	switch err {
	case nil:
		req.FinishCreated("user created")
	case user.ErrLowLoginLength, user.ErrLowPasswordComplexity, user.ErrInvalidRole, processor.ErrUserAlreadyExists:
		req.FinishBadRequest(err.Error())
	case processor.ErrUnexpectedDatabaseBehavior:
		req.FinishError(err.Error())
	default:
		req.FinishError("unexpected error")
	}
}

func (s *Server) handleVerifyUser(w http.ResponseWriter, r *http.Request) {
	req := request.New(w, r)

	id := req.VarsValue("id").MustUInt()
	if id == 0 {
		req.FinishBadRequest("no valid id in request")
		return
	}

	err := s.processor.VerifyUser(r.Context(), id)

	switch err {
	case nil:
		req.FinishOKJSON("user verified")
	case user.ErrAlreadyVerified, processor.ErrUserNotFound:
		req.FinishBadRequest(err.Error())
	case processor.ErrUnexpectedDatabaseBehavior:
		req.FinishError(err.Error())
	default:
		req.FinishError("unexpected error")
	}
}
