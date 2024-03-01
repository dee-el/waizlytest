package v1userhttphandler

import (
	"encoding/json"
	"io"
	"net/http"

	commonerr "waizlytest/common/errors"

	commonhttpenc "waizlytest/common/http/encoder"
	commonhttpresp "waizlytest/common/http/response"

	userservice "waizlytest/services/user"
)

type RegistratorHandler struct {
	registrator userservice.Registrator
}

func NewRegistratorHandler(registrator userservice.Registrator) *RegistratorHandler {
	return &RegistratorHandler{
		registrator: registrator,
	}
}

func (hn *RegistratorHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var p *userservice.RegisterRequest
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

		_, err = hn.registrator.Register(ctx, *p)
		if err != nil {
			commonhttpenc.JSONErrorEncoder(ctx, w, err)
			return
		}

		commonhttpenc.JSONResponseEncoder(ctx, w, http.StatusOK, commonhttpresp.NewResponse(nil, nil))
	}
}
