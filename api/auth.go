package api

import (
	"encoding/json"
	"github.com/marki-eriker/event-listener/pkg/request"
	"github.com/marki-eriker/event-listener/processor"
	"net/http"
)

type loginInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginOutput struct {
	Token string `json:"token"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := request.New(w, r)

	if req.GetBody() == nil || len(req.GetBody()) == 0 {
		req.FinishBadRequest("No JSON body")
		return
	}

	var payload loginInput

	err := json.Unmarshal(req.GetBody(), &payload)
	if err != nil {
		req.FinishBadRequest("Unable to unmarshal JSON data: %s", err)
		return
	}

	token, err := s.processor.LoginUser(r.Context(), payload.Login, payload.Password)

	switch err {
	case nil:
		req.FinishOKJSON(&loginOutput{Token: token})
	case processor.ErrInvalidCredentials:
		req.FinishBadRequest(err.Error())
	case processor.ErrNotVerified:
		req.FinishBadRequest("you must verify your profile. Ask system administrator")
	case processor.ErrUnexpectedDatabaseBehavior, processor.ErrTokenGeneration:
		req.FinishError(err.Error())
	default:
		req.FinishError("unexpected error")
	}
}
