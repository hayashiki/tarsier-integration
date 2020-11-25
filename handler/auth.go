package handler

import (
	"github.com/hayashiki/tarsier-integration/usecase"
	"net/http"
)

func (h *Handler) InvokeSlackAuth(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, h.slackAuthSvc.InvokeURL(), http.StatusTemporaryRedirect)
}

func (h *Handler) HandleSlackAuth(w http.ResponseWriter, r *http.Request) {

	codes, ok := r.URL.Query()["code"]
	if !ok {
		//NewHTTPError(http.StatusBadRequest, "Missing authorization code")
		return
	}

	uc := usecase.NewSlackAuthenticate(h.slackAuthSvc, h.TeamRepo)
	params := usecase.AuthenticateSlackParams{Code: codes[0]}
	resp, err := uc.Do(params)
	if err != nil {
		// NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get slack access err=%w", err))
		return
	}

	w.Write([]byte(h.slackAuthSvc.CallbackHTML(resp.AuthResp.Team.ID)))
	//return jsonResponse(w, http.StatusOK, nil)
}
