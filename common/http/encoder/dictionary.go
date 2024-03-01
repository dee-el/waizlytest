package commonhttpenc

import (
	"net/http"

	commonerr "waizlytest/common/errors"
)

// Dictionary is used to store a mapping from `error.Type` to http status.
// This is to make it easy for user when err happens to its http status as business rule.
type Dictionary map[commonerr.Type]int

var DefaultDictionary Dictionary = Dictionary{
	commonerr.TypeAuthenticationError:   http.StatusUnauthorized,
	commonerr.TypeNotFoundError:         http.StatusNotFound,
	commonerr.TypeForbiddenError:        http.StatusForbidden,
	commonerr.TypeApplicationLimitError: http.StatusTooManyRequests,
	commonerr.TypeInternalServerError:   http.StatusInternalServerError,
	commonerr.TypeMaintenanceError:      http.StatusServiceUnavailable,
	commonerr.TypeBadRequestError:       http.StatusBadRequest,
}

var dictionary = DefaultDictionary

// CreateDictionary is function to create new dictionary.
// This is make user still have flexibiilty to use their own dictionary
func CreateDictionary(d Dictionary) {
	dictionary = d
}
