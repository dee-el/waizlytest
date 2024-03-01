package v1authhttphandler

import (
	"encoding/json"
	"io"
	"net/http"

	commonerr "waizlytest/common/errors"

	commonhttpenc "waizlytest/common/http/encoder"
	commonhttpresp "waizlytest/common/http/response"

	authservice "waizlytest/services/auth"
)

type AuthnHandler struct {
	authn authservice.Authn
}

func NewAuthnHandler(authn authservice.Authn) *AuthnHandler {
	return &AuthnHandler{
		authn: authn,
	}
}

func (hn *AuthnHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var p *authservice.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			if err == io.EOF {
				e := commonerr.New(commonerr.TypeBadRequestError, commonerr.Code101, "request can not be empty")
				commonhttpenc.JSONErrorEncoder(ctx, w, e)
				return
			}

			commonhttpenc.JSONErrorEncoder(ctx, w, err)
			return
		}

		lr, err := hn.authn.Login(ctx, *p)
		if err != nil {
			commonhttpenc.JSONErrorEncoder(ctx, w, err)
			return
		}

		commonhttpenc.JSONResponseEncoder(ctx, w, http.StatusOK, commonhttpresp.NewResponse(lr, nil))
	}
}
