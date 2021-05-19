// Copyright 2020 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth2

import (
	"fmt"
	"strings"
	"time"

	jwtlib "github.com/dgrijalva/jwt-go"
	jwtclaims "github.com/greenpau/caddy-auth-jwt/pkg/claims"
	"github.com/greenpau/caddy-auth-portal/pkg/errors"
)

func (b *Backend) validateAccessToken(state string, data map[string]interface{}) (*jwtclaims.UserClaims, error) {
	var tokenString string
	if v, exists := data[b.IdentityTokenName]; exists {
		tokenString = v.(string)
	} else {
		return nil, fmt.Errorf("token response has no %s field", b.IdentityTokenName)
	}

	token, err := jwtlib.Parse(tokenString, func(token *jwtlib.Token) (interface{}, error) {
		if _, validMethod := token.Method.(*jwtlib.SigningMethodRSA); !validMethod {
			return nil, fmt.Errorf("unsupported signing method: %s", token.Header["alg"])
		}
		keyID, found := token.Header["kid"].(string)
		if !found {
			return nil, fmt.Errorf("kid not found in %s", b.IdentityTokenName)
		}
		key, exists := b.publicKeys[keyID]
		if !exists {
			if !b.disableKeyVerification {
				if err := b.fetchKeysURL(); err != nil {
					return nil, errors.ErrBackendOauthKeyFetchFailed.WithArgs(err)
				}
			}
			key, exists = b.publicKeys[keyID]
			if !exists {
				return nil, fmt.Errorf("the supplied kid not found in jwks public keys: %s", keyID)
			}
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %s", b.IdentityTokenName, err)
	}

	if _, ok := token.Claims.(jwtlib.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenString)
	}

	tokenClaims := token.Claims.(jwtlib.MapClaims)
	if tokenClaims == nil {
		return nil, fmt.Errorf("token claims are nil")
	}

	if _, exists := tokenClaims["nonce"]; !exists {
		return nil, fmt.Errorf("nonce claim not found")
	}
	if err := b.state.validateNonce(state, tokenClaims["nonce"].(string)); err != nil {
		return nil, fmt.Errorf("nonce claim validation failed: %s", err)
	}

	// Create new claims
	claims := &jwtclaims.UserClaims{
		Origin:    b.TokenProvider.TokenOrigin + "/" + state,
		ExpiresAt: time.Now().Add(time.Duration(b.TokenProvider.TokenLifetime) * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Add(10 * time.Minute * -1).Unix(),
	}

	// Iterate over the received token and populate the fields
	flatClaims := []string{
		"sub", "name", "email", "iat", "exp", "jti",
		"iss", "groups", "picture",
	}
	for _, k := range flatClaims {
		if _, exists := tokenClaims[k]; !exists {
			continue
		}
		switch k {
		case "sub":
			claims.Subject = tokenClaims[k].(string)
		case "name":
			claims.Name = tokenClaims[k].(string)
		case "email":
			claims.Email = tokenClaims[k].(string)
		case "iss":
			claims.Origin = tokenClaims[k].(string)
		case "picture":
			claims.PictureURL = tokenClaims[k].(string)
		}
	}

	addRoles := func(i interface{}, path string) error {
		switch roles := i.(type) {
		case []interface{}:
			for _, role := range roles {
				switch role.(type) {
				case string:
					claims.Roles = append(claims.Roles, role.(string))
				default:
					return fmt.Errorf("invalid %s entry type %v", path, role)
				}
			}
		case string:
			for _, role := range strings.Split(roles, " ") {
				claims.Roles = append(claims.Roles, role)
			}
		default:
			return fmt.Errorf("invalid %s type %v", path, i)
		}

		return nil
	}

	nestedClaims := []string{"roles", "role", "groups", "group"}
	for claim, value := range tokenClaims {
		for _, claimName := range nestedClaims {
			if claim == claimName || strings.HasSuffix(claim, "/"+claimName) {
				if err := addRoles(value, claimName); err != nil {
					return nil, err
				}
			}
		}
	}

	if m, exists := tokenClaims["app_metadata"]; exists {
		if metadata, ok := m.(map[string]interface{}); ok {
			for _, claimName := range nestedClaims {
				if _, exists := metadata[claimName]; exists {
					if err := addRoles(metadata[claimName], "app_metadata."+claimName); err != nil {
						return nil, err
					}
				}
			}

			if a, exists := metadata["authorization"]; exists {
				if authorization, ok := a.(map[string]interface{}); ok {
					for _, claimName := range nestedClaims {
						if _, exists := authorization[claimName]; exists {
							if err := addRoles(authorization[claimName], "app_metadata.authorization."+claimName); err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}
	}

	if m, exists := tokenClaims["realm_access"]; exists {
		if metadata, ok := m.(map[string]interface{}); ok {
			for _, claimName := range nestedClaims {
				if _, exists := metadata[claimName]; exists {
					if err := addRoles(metadata[claimName], "realm_access"+claimName); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	if claims.Email == "" {
		return nil, fmt.Errorf("email claim not found")
	}
	if claims.Subject == "" {
		claims.Subject = claims.Email
	}
	if len(claims.Roles) < 1 {
		claims.Roles = []string{"anonymous", "guest", "everyone"}
	}

	return claims, nil
}
