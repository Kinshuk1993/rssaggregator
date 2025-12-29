package auth

import (
	"errors"
	"net/http"
	"strings"
)

// extracts an API Key from headers from a http request
// example:
// Authorization: ApiKey: <api_key_here>
func GetAPIKey(headers *http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authentication info found in request headers")
	}

	vals := strings.Split(authHeader, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header found in request headers")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("no ApiKey key found in request header")
	}

	return vals[1], nil
}
