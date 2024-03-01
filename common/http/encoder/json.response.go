package commonhttpenc

import (
	"context"
	"encoding/json"
	"net/http"

	commonhttpresp "waizlytest/common/http/response"
)

// JSONResponseEncoder encodes the passed response object to the HTTP response writer in JSON format.
func JSONResponseEncoder(ctx context.Context, w http.ResponseWriter, httpStatus int, res *commonhttpresp.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(res)
}
