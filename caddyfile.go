// Copyright 2020 Paul Greenberg @greenpau
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package portal

import (
	"encoding/json"
	"os"
	"strings"

	jwt "github.com/greenpau/caddy-auth-jwt"
	"github.com/greenpau/caddy-auth-portal/pkg/cookies"
	"github.com/greenpau/caddy-auth-portal/pkg/utils"
	ui "github.com/greenpau/caddy-auth-ui"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	httpcaddyfile.RegisterDirective("auth_portal", parseCaddyfileAuthPortal)
}

func initCaddyfileLogger() *zap.Logger {
	logAtom := zap.NewAtomicLevel()
	logAtom.SetLevel(zapcore.DebugLevel)
	logEncoderConfig := zap.NewProductionEncoderConfig()
	logEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logEncoderConfig.TimeKey = "time"
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(logEncoderConfig),
		zapcore.Lock(os.Stdout),
		logAtom,
	))
	return logger

}

// parseCaddyfileAuthPortal sets up an authentication portal. Syntax:
//
//     auth_portal {
//       path /auth
//       context <default|name>
//       backends {
//         local_backend {
//		     method <local>
//		     file <file_path>
//		     realm <name>
//	       }
//	     }
//
//       local_backend <file/path/to/user/db> <realm/name>
//
//	     jwt {
//	       token_name <value>
//	       token_secret <value>
//	       token_issuer <value>
//	     }
//	     ui {
//	       login_template <file_path>
//	       portal_template <file_path>
//	       logo_url <file_path|url_path>
//	       logo_description <value>
//	     }
//
//       cookie_domain <name>
//       cookie_path <name>
//     }
//
func parseCaddyfileAuthPortal(h httpcaddyfile.Helper) ([]httpcaddyfile.ConfigValue, error) {
	portal := AuthPortal{
		PrimaryInstance: true,
		Context:         "default",
		AuthURLPath:     "/auth",
		UserInterface: &UserInterfaceParameters{
			Templates: make(map[string]string),
		},
		Cookies:  &cookies.Cookies{},
		Backends: []Backend{},
	}

	// logger := initCaddyfileLogger()

	for h.Next() {
		args := h.RemainingArgs()
		if len(args) > 0 {
			return nil, h.Errf("auth backend supports only nested args: %v", args)
		}
		for nesting := h.Nesting(); h.NextBlock(nesting); {
			rootDirective := h.Val()
			switch rootDirective {
			case "cookie_domain":
				args := h.RemainingArgs()
				portal.Cookies.Domain = args[0]
			case "cookie_path":
				args := h.RemainingArgs()
				portal.Cookies.Path = args[0]
			case "path":
				args := h.RemainingArgs()
				portal.AuthURLPath = args[0]
			case "context":
				args := h.RemainingArgs()
				portal.Context = args[0]
			case "local_backend":
				args := h.RemainingArgs()
				if len(args) == 0 {
					return nil, h.Errf("auth backend %s directive has no value", rootDirective)
				}
				backendProps := make(map[string]interface{})
				backendProps["method"] = "local"
				backendProps["path"] = args[0]
				if len(args) > 1 {
					backendProps["realm"] = args[1]
				} else {
					backendProps["realm"] = "local"
				}
				backendJSON, err := json.Marshal(backendProps)
				if err != nil {
					return nil, h.Errf("auth backend %s directive failed to compile to JSON: %s", rootDirective, err.Error())
				}
				backend := Backend{}
				if err := backend.UnmarshalJSON(backendJSON); err != nil {
					return nil, h.Errf("auth backend %s directive failed to compile to JSON: %s", rootDirective, err.Error())
				}
				portal.Backends = append(portal.Backends, backend)
			case "backends":
				for nesting := h.Nesting(); h.NextBlock(nesting); {
					backendName := h.Val()
					backendProps := make(map[string]interface{})
					var backendAuthMethod string
					for subNesting := h.Nesting(); h.NextBlock(subNesting); {
						backendArg := h.Val()
						switch backendArg {
						case "method", "type":
							if !h.NextArg() {
								return nil, h.Errf("auth backend %s subdirective %s has no value", backendName, backendArg)
							}
							backendAuthMethod = h.Val()
							backendProps["method"] = backendAuthMethod
						case "username", "password", "search_base_dn", "search_filter", "path", "realm":
							if !h.NextArg() {
								return nil, h.Errf("auth backend %s subdirective %s has no value", backendName, backendArg)
							}
							backendProps[backendArg] = h.Val()
						case "attributes":
							attrMap := make(map[string]interface{})
							for attrNesting := h.Nesting(); h.NextBlock(attrNesting); {
								attrName := h.Val()
								if !h.NextArg() {
									return nil, h.Errf("auth backend %s subdirective %s key %s has no value", backendName, backendArg, attrName)
								}
								attrMap[attrName] = h.Val()
							}
							backendProps[backendArg] = attrMap
						case "servers":
							serverMaps := []map[string]interface{}{}
							for serverNesting := h.Nesting(); h.NextBlock(serverNesting); {
								serverMap := make(map[string]interface{})
								serverMap["addr"] = h.Val()
								serverProps := h.RemainingArgs()
								if len(serverProps) > 0 {
									for _, serverProp := range serverProps {
										switch serverProp {
										case "ignore_cert_errors":
											serverMap[serverProp] = true
										default:
											return nil, h.Errf("auth backend %s subdirective %s prop %s is unsupported", backendName, backendArg, serverProp)
										}
									}
								}
								serverMaps = append(serverMaps, serverMap)
							}
							backendProps[backendArg] = serverMaps
						case "groups":
							groupMaps := []map[string]interface{}{}
							for groupNesting := h.Nesting(); h.NextBlock(groupNesting); {
								groupMap := make(map[string]interface{})
								groupDN := h.Val()
								groupMap["dn"] = groupDN
								groupRoles := h.RemainingArgs()
								if len(groupRoles) == 0 {
									return nil, h.Errf("auth backend %s subdirective %s dn %s has no roles", backendName, backendArg, groupDN)
								}
								groupMap["roles"] = groupRoles
								groupMaps = append(groupMaps, groupMap)
							}
							backendProps[backendArg] = groupMaps
						default:
							return nil, h.Errf("unknown auth backend %s subdirective: %s", backendName, backendArg)
						}
					}
					backendJSON, err := json.Marshal(backendProps)
					if err != nil {
						return nil, h.Errf("auth backend %s subdirective failed to compile to JSON: %s", backendName, err.Error())
					}
					backend, err := NewBackendFromBytes(backendAuthMethod, backendJSON)
					if err != nil {
						return nil, h.Errf("auth backend %s subdirective failed to compile backend from JSON: %s", backendName, err.Error())
					}
					portal.Backends = append(portal.Backends, *backend)
				}
			case "jwt":
				if portal.TokenProvider == nil {
					portal.TokenProvider = &jwt.TokenProviderConfig{}
				}
				for nesting := h.Nesting(); h.NextBlock(nesting); {
					subDirective := h.Val()
					switch subDirective {
					case "token_name":
						if !h.NextArg() {
							return nil, h.Errf("%s %s subdirective has no value", rootDirective, subDirective)
						}
						portal.TokenProvider.TokenName = h.Val()
					case "token_secret":
						if !h.NextArg() {
							return nil, h.Errf("%s %s subdirective has no value", rootDirective, subDirective)
						}
						portal.TokenProvider.TokenSecret = h.Val()
					case "token_issuer":
						if !h.NextArg() {
							return nil, h.Errf("%s %s subdirective has no value", rootDirective, subDirective)
						}
						portal.TokenProvider.TokenIssuer = h.Val()
					default:
						return nil, h.Errf("unknown subdirective for %s: %s", rootDirective, subDirective)
					}
				}
			case "ui":
				for nesting := h.Nesting(); h.NextBlock(nesting); {
					subDirective := h.Val()
					if strings.HasSuffix(subDirective, "_template") {
						if !h.NextArg() {
							return nil, h.Errf("%s %s subdirective has no value", rootDirective, subDirective)
						}
						templateName := strings.TrimSuffix(subDirective, "_template")
						portal.UserInterface.Templates[templateName] = h.Val()
					} else {
						switch subDirective {
						case "logo_url":
							if !h.NextArg() {
								return nil, h.Errf("%s %s subdirective has no value", rootDirective, subDirective)
							}
							portal.UserInterface.LogoURL = h.Val()
						case "logo_description":
							if !h.NextArg() {
								return nil, h.Errf("%s %s subdirective has no value", rootDirective, subDirective)
							}
							portal.UserInterface.LogoDescription = h.Val()
						case "auto_redirect_url":
							if !h.NextArg() {
								return nil, h.Errf("%s %s subdirective has no value", rootDirective, subDirective)
							}
							portal.UserInterface.AutoRedirectURL = h.Val()
						case "links":
							for subNesting := h.Nesting(); h.NextBlock(subNesting); {
								title := h.Val()
								args := h.RemainingArgs()
								if len(args) == 0 {
									return nil, h.Errf("auth backend %s subdirective %s has no value", subDirective, title)
								}
								portal.UserInterface.PrivateLinks = append(portal.UserInterface.PrivateLinks, ui.UserInterfaceLink{
									Title: title,
									Link:  args[0],
								})
							}
						default:
							return nil, h.Errf("unsupported subdirective for %s: %s", rootDirective, subDirective)
						}
					}
				}
			default:
				return nil, h.Errf("unsupported root directive: %s", rootDirective)
			}
		}
	}

	if portal.AuthURLPath == "" {
		portal.AuthURLPath = "/auth"
	}
	if strings.HasSuffix(portal.AuthURLPath, "*") {
		return nil, h.Errf("path directive must not end with '*', got %s", portal.AuthURLPath)
	}
	if !strings.HasPrefix(portal.AuthURLPath, "/") {
		return nil, h.Errf("path directive must begin with '/', got %s", portal.AuthURLPath)
	}

	if portal.Context == "" {
		return nil, h.Errf("context directive must not be empty")
	}

	if portal.TokenProvider == nil {
		portal.TokenProvider = &jwt.TokenProviderConfig{}
		portal.TokenProvider.TokenSecret = utils.GetRandomStringFromRange(32, 64)
	}

	h.Reset()
	h.Next()
	pathMatcher := caddy.ModuleMap{
		"path": h.JSON(caddyhttp.MatchPath{portal.AuthURLPath + "*"}),
	}
	route := caddyhttp.Route{
		HandlersRaw: []json.RawMessage{caddyconfig.JSONModuleObject(portal, "handler", "auth_portal", nil)},
	}
	subroute := new(caddyhttp.Subroute)
	subroute.Routes = append([]caddyhttp.Route{route}, subroute.Routes...)
	return h.NewRoute(pathMatcher, subroute), nil
}
