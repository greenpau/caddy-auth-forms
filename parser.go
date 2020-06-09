package forms

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func parseCredentials(r *http.Request) (map[string]string, error) {
	if r.Method == "POST" {
		return parseAuthForm(r)
	}
	if r.Method == "GET" {
		return parseAuthRequest(r)
	}
	return nil, fmt.Errorf("Request method %s is unsupported", r.Method)
}

func parseAuthForm(r *http.Request) (map[string]string, error) {
	var reqFields []string
	kv := make(map[string]string)
	var maxBytesLimit int64 = 1000
	var minBytesLimit int64 = 15
	if r.ContentLength > maxBytesLimit {
		return nil, fmt.Errorf("Request payload exceeded the limit of %d bytes: %d", maxBytesLimit, r.ContentLength)
	}
	if r.ContentLength < minBytesLimit {
		return nil, fmt.Errorf("Request payload is too small: %d", r.ContentLength)
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/x-www-form-urlencoded" {
		return nil, fmt.Errorf("Request content type is not application/x-www-form-urlencoded")
	}

	rq := r.FormValue("activity")
	if rq == "" {
		rq = "login"
	}

	switch rq {
	case "login":
		reqFields = []string{"username", "password", "realm"}
	default:
		return nil, fmt.Errorf("request type is unsupported")
	}

	for _, k := range reqFields {
		if v := r.FormValue(k); v != "" {
			kv[k] = v
		}
	}

	if _, exists := kv["realm"]; !exists {
		kv["realm"] = "local"
	}

	return kv, nil
}

func parseAuthRequest(r *http.Request) (map[string]string, error) {
	kv := make(map[string]string)
	authzHeaderStr := r.Header.Get("Authorization")
	if authzHeaderStr == "" {
		return nil, nil
	}

	authzHeaderParts := strings.Split(authzHeaderStr, ",")
	if len(authzHeaderParts) == 0 {
		return nil, nil
	}

	authzStrParts := strings.Split(authzHeaderParts[0], " ")
	if len(authzStrParts) != 2 {
		return nil, nil
	}

	authzType := authzStrParts[0]
	if authzType != "Basic" {
		return nil, nil
	}
	authzStr, err := base64.StdEncoding.DecodeString(authzStrParts[1])
	if err != nil {
		return nil, fmt.Errorf("decoding error: %s", err)
	}
	authzArr := strings.SplitN(string(authzStr), ":", 2)
	if len(authzArr) != 2 {
		return nil, fmt.Errorf("parsing error: %s", err)
	}
	kv["username"] = authzArr[0]
	kv["password"] = authzArr[1]
	if len(authzHeaderParts) == 1 {
		kv["realm"] = "local"
		return kv, nil
	}
	realmHeaderParts := strings.Split(authzHeaderParts[1], "=")
	if len(realmHeaderParts) != 2 {
		return nil, fmt.Errorf("realm parsing failed for %s", realmHeaderParts)
	}
	if realmHeaderParts[0] != "realm" {
		return nil, fmt.Errorf("realm not found in %s", realmHeaderParts)
	}
	kv["realm"] = realmHeaderParts[1]
	return kv, nil
}
