package commonhttpenc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	commonerr "waizlytest/common/errors"
	commonhttpresp "waizlytest/common/http/response"
)

// JSONResponseEncoder encodes the passed err to client in JSON format.
// Using Dictionary to help directing err to each own HTTP status.
func JSONErrorEncoder(ctx context.Context, w http.ResponseWriter, err error) {
	// empty response, value will fill from type checking err
	// same response to standardize format response
	resp := commonhttpresp.NewResponse(nil, nil)

	httpStatus := http.StatusOK
	if err != nil {
		switch e := err.(type) {
		case *commonerr.Error:
			resp.Error = e

			httpStatus = dictionary[e.Type]
		default:

			log.Printf("[InternalServerError]: %v", err)
			// mask the real error
			// client(s) doesn't need to know what it is
			resp.Error = commonerr.ErrorInternalServer
			httpStatus = dictionary[resp.Error.Type]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	// no need check err encoder
	json.NewEncoder(w).Encode(resp)
}
