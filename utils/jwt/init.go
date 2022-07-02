package jwt

import (
	"errors"
	"strings"
)

func DecodeBearer(authorization string) (string, error) {
	if len(authorization) == 0 {
		return "", errors.New("missing authorization header")
	}
	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("malformed authorization header")
	}
	token := strings.Replace(authorization, "Bearer ", "", 1)
	if len(token) == 0 {
		return "", errors.New("malformed authorization header")
	}
	return token, nil
}
