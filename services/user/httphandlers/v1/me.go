package v1userhttphandler

import (
	"encoding/json"
	"io"
	"net/http"

	"waizlytest/common/contextkey"
	commonerr "waizlytest/common/errors"

	commonhttpenc "waizlytest/common/http/encoder"
	commonhttpresp "waizlytest/common/http/response"

	userservice "waizlytest/services/user"
)

type MeHandler struct {
	me userservice.Me
}

func NewMeHandler(me userservice.Me) *MeHandler {
	return &MeHandler{
		me: me,
	}
}

func (hn *MeHandler) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// this cast is safe, as we controlled the auth middleware
		id, ok := ctx.Value(contextkey.UserIDContextKey).(int64)
		if !ok {
			e := commonerr.ErrorForbidden
			commonhttpenc.JSONErrorEncoder(ctx, w, e)
			return
		}

		pr, err := hn.me.GetProfile(ctx, id)
		if err != nil {
			commonhttpenc.JSONErrorEncoder(ctx, w, err)
			return
		}

		commonhttpenc.JSONResponseEncoder(ctx, w, http.StatusOK, commonhttpresp.NewResponse(pr, nil))
	}
}

func (hn *MeHandler) UpdateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// this cast is safe, as we controlled the auth middleware
		id, ok := ctx.Value(contextkey.UserIDContextKey).(int64)
		if !ok {
			e := commonerr.ErrorForbidden
			commonhttpenc.JSONErrorEncoder(ctx, w, e)
			return
		}

		var p *userservice.UpdateRequest
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

		err = hn.me.UpdateProfile(ctx, id, *p)
		if err != nil {
			commonhttpenc.JSONErrorEncoder(ctx, w, err)
			return
		}

		commonhttpenc.JSONResponseEncoder(ctx, w, http.StatusOK, commonhttpresp.NewResponse(nil, nil))
	}
}
