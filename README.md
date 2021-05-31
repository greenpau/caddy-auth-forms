# caddy-auth-portal

<a href="https://github.com/greenpau/caddy-auth-portal/actions/" target="_blank"><img src="https://github.com/greenpau/caddy-auth-portal/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/caddy-auth-portal" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>
<a href="https://caddy.community" target="_blank"><img src="https://img.shields.io/badge/community-forum-ff69b4.svg"></a>

Authentication Plugin for [Caddy v2](https://github.com/caddyserver/caddy) implementing
Form-Based, Basic, Local, LDAP, OpenID Connect, OAuth 2.0, SAML Authentication.

Please see other relevant plugins:
* [caddy-auth-jwt](https://github.com/greenpau/caddy-auth-jwt)
* [caddy-trace](https://github.com/greenpau/caddy-trace)

Download Caddy with the plugins enabled:

* <a href="https://caddyserver.com/api/download?os=linux&arch=amd64&p=github.com%2Fgreenpau%2Fcaddy-auth-jwt%40v1.2.7&p=github.com%2Fgreenpau%2Fcaddy-auth-portal%40v1.4.5&p=github.com%2Fgreenpau%2Fcaddy-trace%40v1.1.6" target="_blank">linux/amd64</a>
* <a href="https://caddyserver.com/api/download?os=windows&arch=amd64&p=github.com%2Fgreenpau%2Fcaddy-auth-jwt%40v1.2.7&p=github.com%2Fgreenpau%2Fcaddy-auth-portal%40v1.4.5&p=github.com%2Fgreenpau%2Fcaddy-trace%40v1.1.6" target="_blank">windows/amd64</a>

Please show your appreciation for this work and :star: :star: :star:

Please ask questions either here or via LinkedIn. I am happy to help you! @greenpau

<!-- begin-markdown-toc -->
## Table of Contents

* [Overview](#overview)
* [Authentication Portal](#authentication-portal)
  * [User Identity](#user-identity)
  * [User Settings](#user-settings)
  * [Multi-Factor Authentication MFA](#multi-factor-authentication-mfa)
    * [Enabling MFA](#enabling-mfa)
    * [Add MFA Authenticator Application](#add-mfa-authenticator-application)
  * [Theming](#theming)
* [Authorization Cookie](#authorization-cookie)
  * [Intra-Domain Cookies](#intra-domain-cookies)
  * [JWT Tokens](#jwt-tokens)
    * [JWT Signing Method](#jwt-signing-method)
* [Usage Examples](#usage-examples)
  * [Secure Prometheus](#secure-prometheus)
  * [Secure Kibana](#secure-kibana)
* [Authentication Methods](#authentication-methods)
  * [Basic Authentication](#basic-authentication)
  * [Form-Based Authentication](#form-based-authentication)
* [User Interface Features](#user-interface-features)
  * [Auto-Redirect URL](#auto-redirect-url)
  * [User Registration](#user-registration)
  * [Custom CSS Styles](#custom-css-styles)
  * [Custom Javascript](#custom-javascript)
  * [Portal Links](#portal-links)
  * [Custom Header](#custom-header)
  * [Static Assets of Any Type](#static-assets-of-any-type)
* [Local Authentication Backend](#local-authentication-backend)
  * [Configuration Primer](#configuration-primer)
  * [Identity Store](#identity-store)
  * [Password Management](#password-management)
* [LDAP Authentication Backend](#ldap-authentication-backend)
  * [Configuration Primer](#configuration-primer-1)
  * [LDAP Authentication Process](#ldap-authentication-process)
* [SAML Authentication Backend](#saml-authentication-backend)
  * [Time Synchronization](#time-synchronization)
  * [Configuration](#configuration)
  * [Azure Active Directory (Office 365) Applications](#azure-active-directory-office-365-applications)
    * [Azure AD SAML Configuration](#azure-ad-saml-configuration)
    * [Set Up Azure AD Application](#set-up-azure-ad-application)
    * [Configure SAML Authentication](#configure-saml-authentication)
    * [Azure AD IdP Metadata and Certificate](#azure-ad-idp-metadata-and-certificate)
    * [User Interface Options](#user-interface-options)
    * [Development Notes](#development-notes)
* [OAuth 2.0 and OpenID Connect (OIDC) Authentication Backend](#oauth-20-and-openid-connect-oidc-authentication-backend)
  * [OAuth 2.0 Flow](#oauth-20-flow)
  * [Adding Role Claims](#adding-role-claims)
  * [OAuth 2.0 Authorization Servers and Identity Providers](#oauth-20-authorization-servers-and-identity-providers)
    * [Okta](#okta)
    * [Google Identity Platform](#google-identity-platform)
    * [LinkedIn](#linkedin)
    * [Auth0](#auth0)
    * [OneLogin](#onelogin)
    * [Microsoft](#microsoft)
    * [Github](#github)
    * [Facebook](#facebook)
* [X.509 Certificate-based Authentication Backend](#x509-certificate-based-authentication-backend)
* [Miscellaneous](#miscellaneous)
  * [Binding to Privileged Ports](#binding-to-privileged-ports)
  * [Recording Source IP Address in JWT Token](#recording-source-ip-address-in-jwt-token)
  * [Session ID Cache](#session-id-cache)
  * [Caddyfile Shortcuts](#caddyfile-shortcuts)
  * [Multiple Hosts](#multiple-hosts)

<!-- end-markdown-toc -->

## Overview

The purpose of this plugin is providing **authentication** only. The plugin
issue JWT tokens upon successful authentication. In turn, the **authorization**
of the tokens is being handled by [`caddy-auth-jwt`](https://github.com/greenpau/caddy-auth-jwt).

The plugin supports the following **authentication** backends:

* Local (`local`) - JSON flat file database
* LDAP (`ldap`) - remote Microsoft AD database
* SAML
* OAuth 2.0 and OpenID Connect (OIDC)

The plugin accepts user credentials for **authentication** with:

* Form-based Authentication: `POST` with `application/x-www-form-urlencoded`
* Basic Authentication: `GET` with `Authorization: Basic` header

The following digram is visual representation of the configuration of
[`caddy-auth-portal`](https://github.com/greenpau/caddy-auth-portal) and
[`caddy-auth-jwt`](https://github.com/greenpau/caddy-auth-jwt).

![Authentication Plugins](assets/docs/images/auth_plugin_arch.png)

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->

## Authentication Portal

<table cellspacing="0" cellpadding="0" style="border-collapse: collapse; border: none; vertical-align: top;">
  <tr style="border: none;">
    <td style="border: none; vertical-align: top;">
      <img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/basic_login.png">
    </td>
    <td style="border: none; vertical-align: top;">
      <img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/basic_portal.png">
    </td>
  </tr>
</table>

### User Identity

The following screenshot is from `/auth/whoami` endpoint:

<img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/whoami.png">

### User Settings

The following screenshot is from `/auth/settings/` endpoint:

<img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/settings.png">

### Multi-Factor Authentication MFA

#### Enabling MFA

MFA can be enabled globally and on per-backend basis via `require mfa` directive.

```
        local_backend {
          method local
          path /etc/gatekeeper/auth/local/users_db.json
          realm local
          require mfa
        }
```

Currently, the MFA requirement can be applied only to `local` backend type.

Importantly, a user may choose to enable MFA via the user's settings page.

#### Add MFA Authenticator Application

The following screenshot is from `/auth/settings/mfa/add/app` endpoint:

<img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/settings_mfa_app.png">

The QR Code displayed on the page complies [Key Uri Format](https://github.com/google/google-authenticator/wiki/Key-Uri-Format).

In your MFA application, e.g. Microsoft Authenticator, follow these steps to
onboard your web account.

<table cellspacing="0" cellpadding="0" style="border-collapse: collapse; border: none; vertical-align: top;">
  <tr style="border: none;">
    <td style="border: none; vertical-align: top;">
      <img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/ms_mfa_app_add_account.png">
    </td>
    <td style="border: none; vertical-align: top;">
      <img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/ms_mfa_app_scan_qrcode.png">
    </td>
    <td style="border: none; vertical-align: top;">
      <img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/ms_mfa_app_new_account.png">
    </td>
  </tr>
</table>

### Theming

The theming of the portal works as follows.

It starts with a concept of `theme`. By default, the portal uses `basic` theme.
There is no need to defind it in Caddyfile.

```
localhost {
  route /auth* {
    auth_portal {
      ui {
        theme basic
      }
```

Each theme must have a set of default pages:

* `generic`
* `login`
* `portal`
* `register`
* `whoami`
* `settings`
* `sandbox`

The plain text templates are being stored in `assets/templates/<THEME>/<PAGE>.template`.

```
assets/templates/basic/generic.template
assets/templates/basic/login.template
assets/templates/basic/portal.template
assets/templates/basic/register.template
assets/templates/basic/whoami.template
assets/templates/basic/settings.template
assets/templates/basic/sandbox.template
```

These templates are the parts of `pkg/ui/pages.go`. They are compiled in the
portal's binary. That is, there is no need to store them on the disk.

Next, if a user wants to use a different template, then it could be passed via
Caddyfile directives. Specifically, use `<PAGE>_template` directive to point
to a file on disk.

```
localhost {
  route /auth* {
    auth_portal {
      ui {
        theme basic
        login_template "/etc/gatekeeper/ui/login.template"
      }
```

TODO: Review [Refactoring UI Feed](https://twitter.com/i/events/994601867987619840)
and [Refactoring UI Website](https://refactoringui.com/).

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->

## Authorization Cookie

### Intra-Domain Cookies

The following `Caddyfile` settings define the scope of the cookies issued by
the plugin. Specifically, what URLs the cookies should be sent to.
See [MDN - Using HTTP cookies - Define where cookies are sent](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies)
for more information.


* `cookie_domain`: adds the **Domain** attribute to a cookie. It determines
  which hosts are allowed to receive the cookie.
* `cookie_path`: adds the **Path** attribute to a cookie. It determines the
  URL path that must exist in the requested URL in order to send
  the Cookie header.
* `cookie_lifetime`: sets the number of seconds until the cookie expires.
  The directive sets "Max-Age" cookie attribute.

### JWT Tokens

The plugin sends JWT token via the cookie.

* `token_name`: specifies the names of the cookie with authorization credentials

By default the lifetime of the token is 15 minutes. The `token_lifetime`
can be used to change it to 1 hour (3600 seconds).

```
      jwt {
        token_name access_token
        token_secret 0e2fdcf8-6868-41a7-884b-7308795fc286
        token_lifetime 3600
      }
```

The issued JWT token could be of two types:

1. `HS512`: signed using shared secret key
2. `RS512`: signed using private PEM key

The `HS512` is being configured with `token_secret`

```
      jwt {
        ...
        # token_secret <shared_key>
        token_secret 0e2fdcf8-6868-41a7-884b-7308795fc286
        ...
      }
```

The `RS512` is being configured with `token_rsa_file` directive:

```
      jwt {
        ...
        token_rsa_file <key_id> <file_path>
        token_rsa_file Hz789bc303f0db /etc/gatekeeper/auth/jwt/sign_key.pem
        ...
      }
```

If necessary, generate the signing key:

```bash
$ openssl genrsa -out /etc/gatekeeper/auth/jwt/sign_key.pem 2048
Generating RSA private key, 2048 bit long modulus (2 primes)
.....................................................................................................................+++++
....+++++
e is 65537 (0x010001)
```

#### JWT Signing Method

By default, the plugin uses HS512 (shared secret) and RS512 (public/private keys) for
the signing of JWT tokens. User `token_sign_method` to change the algorithm, e.g.

```
      jwt {
        ...
        token_secret 0e2fdcf8-6868-41a7-884b-7308795fc286
        token_sign_method HS256
        ...
      }
```

or:

```
      jwt {
        ...
        token_rsa_file Hz789bc303f0db /etc/gatekeeper/auth/jwt/sign_key.pem
        token_sign_method RS256
        ...
      }
```

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->

## Usage Examples

### Secure Prometheus

The following `Caddyfile` secures Prometheus/Alertmanager services:

```Caddyfile
{
  http_port     8080
  https_port    8443
  debug
}

localhost:8443 {
  route /auth* {
    auth_portal {
      path /auth
      backends {
        local_backend {
          method local
          path /etc/gatekeeper/auth/local/users.json
          realm local
        }
      }
      jwt {
        token_name access_token
        token_secret 0e2fdcf8-6868-41a7-884b-7308795fc286
      }
      ui {
        links {
          "Prometheus" /prometheus
          "Alertmanager" /alertmanager
          "My App" /myapp
        }
      }
    }
  }

  route /prometheus* {
    jwt {
      primary yes
      trusted_tokens {
        static_secret {
          token_name access_token
          token_secret 0e2fdcf8-6868-41a7-884b-7308795fc286
        }
      }
      auth_url /auth
      allow roles anonymous guest admin
      allow roles superadmin
    }
    uri strip_prefix /prometheus
    reverse_proxy http://127.0.0.1:9080
  }

  route /alertmanager* {
    jwt
    uri strip_prefix /alertmanager
    reverse_proxy http://127.0.0.1:9083
  }

  route /myapp* {
    jwt
    respond * "myapp" 200
  }

  route /version* {
    respond * "1.0.0" 200
  }

  route {
    redir https://{hostport}/auth 302
  }
}
```

If you would like to style the UI differently, then specify your
templates and settings:

```
      ui {
        login_template "/etc/gatekeeper/ui/login.template"
        portal_template "/etc/gatekeeper/ui/portal.template"
        logo_url "https://caddyserver.com/resources/images/caddy-circle-lock.svg"
        logo_description "Caddy"
        links {
          "Prometheus" /prometheus
          "Alertmanager" /alertmanager
          "My App" /myapp
        }
      }
```

In fact, if you are not going to display any links, then
remove the `ui` section and use an auto-redirect feature.

[:arrow_up: Back to Top](#table-of-contents)

### Secure Kibana

First, add the following line in `/etc/kibana/kibana.yml`. It must match the
the prefix used when proxying traffic through:

```
server.basePath: "/elk"
```

Next, add the following route in you Caddyfile:

```
  route /elk* {
    jwt
    uri strip_prefix /elk
    reverse_proxy KIBANA_IP:5601
  }
```

Also, add the link to Kibana in `ui` section of Caddyfile:

```
      ui {
        ...
        links {
          ...
          "Kibana" /elk/
          ...
        }
      }
```

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->

## Authentication Methods

### Basic Authentication

The following command demonstrates basic authentication process.
The plugin returns JWT token via `Set-Cookie: access_token` and
`token` field in JSON response.

```bash
curl --insecure -H "Accept: application/json" --user webadmin:password123 -v https://127.0.0.1:3443/auth
```

The expected output is as follows:

```
* About to connect() to 127.0.0.1 port 3443 (#0)
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 3443 (#0)
* Initializing NSS with certpath: sql:/etc/pki/nssdb
* skipping SSL peer certificate verification
* SSL connection using TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
* Server certificate:
*       subject: E=admin@caddy.local,OU=Local Developement,CN=*.caddy.localhost,L=Local Developement,O=Local Developement,ST=NY,C=US
*       start date: Mar 02 08:01:16 2020 GMT
*       expire date: Feb 28 08:01:16 2030 GMT
*       common name: *.caddy.localhost
*       issuer: E=admin@caddy.local,OU=Local Developement,CN=*.caddy.localhost,L=Local Developement,O=Local Developement,ST=NY,C=US
* Server auth using Basic with user 'webadmin'
> GET /auth HTTP/1.1
> Authorization: Basic d2ViYWRtaW46cGFzc3dvcmQxMjM=
> User-Agent: curl/7.29.0
> Host: 127.0.0.1:3443
> Accept: application/json
>
< HTTP/1.1 200 OK
< Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE3MzE0NzksInN1YiI6IndlYmFkbWluIiwiZW1haWwiOiJ3ZWJhZG1pbkBsb2NhbGRvbWFpbi5sb2NhbCIsInJvbGVzIjpbInN1cGVyYWRtaW4iLCJndWVzdCIsImFub255bW91cyJdLCJvcmlnaW4iOiJsb2NhbGhvc3QifQ.OmFOCu-UJdx16FYLa2ezr7WRmOdUbgrQadhfk1tN4AliIwu69x9TLgzoke_Cr3TqzvMjlQDd22r-3DHBXuzllw
< Cache-Control: no-store
< Content-Type: application/json
< Pragma: no-cache
< Server: Caddy
< Set-Cookie: access_token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE3MzE0NzksInN1YiI6IndlYmFkbWluIiwiZW1haWwiOiJ3ZWJhZG1pbkBsb2NhbGRvbWFpbi5sb2NhbCIsInJvbGVzIjpbInN1cGVyYWRtaW4iLCJndWVzdCIsImFub255bW91cyJdLCJvcmlnaW4iOiJsb2NhbGhvc3QifQ.OmFOCu-UJdx16FYLa2ezr7WRmOdUbgrQadhfk1tN4AliIwu69x9TLgzoke_Cr3TqzvMjlQDd22r-3DHBXuzllw Secure; HttpOnly;
< Date: Tue, 09 Jun 2020 19:22:59 GMT
< Content-Length: 318
<
* Connection #0 to host 127.0.0.1 left intact
{"token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE3MzE0NzksInN1YiI6IndlYmFkbWluIiwiZW1haWwiOiJ3ZWJhZG1pbkBsb2NhbGRvbWFpbi5sb2NhbCIsInJvbGVzIjpbInN1cGVyYWRtaW4iLCJndWVzdCIsImFub255bW91cyJdLCJvcmlnaW4iOiJsb2NhbGhvc3QifQ.OmFOCu-UJdx16FYLa2ezr7WRmOdUbgrQadhfk1tN4AliIwu69x9TLgzoke_Cr3TqzvMjlQDd22r-3DHBXuzllw"}
```

### Form-Based Authentication

TBD.

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->
## User Interface Features

### Auto-Redirect URL

Consider the following configuration snippet. When the JWT plugin detects
unauthenticated user, it forwards the user to `https://auth.example.com`.
The `redirect_url` in URL query creates `AUTH_PORTAL_REDIRECT_URL` cookie
in the users session. Upon successful authentication, the portal
clears the cookie and redirects the user to the path specified in
`AUTH_PORTAL_REDIRECT_URL` cookie.

```
https://chat.example.com {
  jwt {
    auth_url https://auth.example.com/auth?redirect_url=https://chat.example.com
  }
}
```

[:arrow_up: Back to Top](#table-of-contents)

### User Registration

The following Caddy configuration enables user registration.

```
registration {
  dropbox /etc/gatekeeper/auth/local/registrations_db.json
  title "User Registration"
  code "NY2020"
  require accept_terms
  require domain_mx
}
```

The parameters are:

* `dropbox`: The file path pointing to registration database.
* `code`: The registration code. A user must know what that code is to
  successfully submit a registration request.
* `require accept_terms`: A user must accept terms and conditions, as well
  as privacy policy to proceed
* `disabled on`: disables user registration
* `title`: changes the title of the registration page
* `require domain_mx`: forces the check of domain MX record

This screenshot is the registration screen with default options:

<img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/portal_registration_simple.png">

The following is the registration screen with mandatory registration
code and the acceptable of terms and conditions:

<img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/portal_registration_terms_code.png">

[:arrow_up: Back to Top](#table-of-contents)

### Custom CSS Styles

The following Caddyfile directive adds a custom CSS stylesheet to the
plugin's pages. The stylesheet is available under `auth/assets/css/custom.css`

```
      ui {
        ...
        custom_css_path path/to/styles.css
        ...
      }
```

[:arrow_up: Back to Top](#table-of-contents)

### Custom Javascript

The following Caddyfile directive adds a custom javascript file to the
plugin's pages. The script is available under `auth/assets/js/custom.js`

```
      ui {
        ...
        custom_js_path path/to/script.js
        ...
      }
```

This directive is usefule for adding Google Analytics or other
minor javascript code.

[:arrow_up: Back to Top](#table-of-contents)

### Portal Links

The following Caddyfile directive sets links that a user would see
upon a successful login:

```bash
      ui {
        ...
        links {
          "Prometheus" /prometheus
          "Alertmanager" /alertmanager
          "My App" /myapp
        }
        ...
      }
```

The link can be opened in a new tab or window via `target_blank` argument:

```
          "My App" /myapp target_blank
```

The link can be disabled with `disabled` argument:

```
          "My App" /myapp disabled
```

The link can have an icon associated with it via `icon` argument:

```
          "My App" /myapp icon "las la-cog"
```

The icon is the reference to [Line Awesome](https://icons8.com/line-awesome) by Icon8.

![Portal - UI - Icons](./assets/docs/images/portal_ui_icons.png)

[:arrow_up: Back to Top](#table-of-contents)

### Custom Header

The following Caddyfile directive injects the code found in `path/to/head.html`
to `<head>` section of the portal's pages:

```bash
      ui {
        ...
        custom_html_header_path path/to/head.html
        ...
      }
```

[:arrow_up: Back to Top](#table-of-contents)

### Static Assets of Any Type

The following Caddyfile directive is capable of loading and serving any type of static
asset, e.g. `js`, `css`, etc.

```bash
      ui {
        ...
        static_asset "assets/css/app.css" "text/css" /path/to/app/styles.css
        ...
      }
```

The above configuration would cause the plugin to read `/path/to/app/styles.css`
and begin serving it with content type of `text/css`  at
`AUTH_PORTAL/assets/css/app.css`, e.g. `https://localhost:8443/auth/assets/css/app.css`.

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->

## Local Authentication Backend

### Configuration Primer

Please refer to the `assets/conf/local/config.json` configuration file when
configurin the plugin backend. In the example, the route refers to `local` backend in
the file `assets/backends/local/users.json`. Specify the path to the file
where you want your database to reside. Do not create a file, but rather
create leading directories.

For example, create `/etc/caddy/auth/local` directory and specify the
`path` value as:

```json
"path": "/etc/caddy/auth/local/users.json",
```

Next, start the server, and find the following following log entries:

```json
{"level":"info","ts":1588704471.5784082,"logger":"http.authentication.providers.portal","msg":"created new user","user_id":"cd5f647a-cc04-4ae2-9d0a-2d5e9b95cf98","user_name":"webadmin","user_email":"webadmin@localdomain.local","user_claims":{"roles":"superadmin"}}
{"level":"info","ts":1588704471.5784378,"logger":"http.authentication.providers.portal","msg":"created default superadmin user for the database","user_name":"webadmin","user_secret":"d87e7749-0dd8-482b-91a2-ada370263293"}
```

### Identity Store

The `user_name` and `user_secret` are password for the `superuser` in the database.

The plugin creates the following a file having the following structure.

```json
{
  "revision": 1,
  "users": [
    {
      "id": "cd5f647a-cc04-4ae2-9d0a-2d5e9b95cf98",
      "username": "webadmin",
      "email_addresses": [
        {
          "address": "webadmin@localdomain.local",
          "domain": "localdomain.local"
        }
      ],
      "passwords": [
        {
          "purpose": "generic",
          "type": "bcrypt",
          "hash": "$2a$10$B67nHY0PEdxLYdyoLk1YLOomvs.T/dSIyzPuoX9vWULrsD3PRf/sq",
          "cost": 10,
          "expired_at": "0001-01-01T00:00:00Z",
          "created_at": "2020-05-05T18:47:51.513552501Z",
          "disabled_at": "0001-01-01T00:00:00Z"
        }
      ],
      "created": "2020-05-05T18:47:51.513552066Z",
      "last_modified": "2020-05-05T18:47:51.513552175Z",
      "roles": [
        {
          "name": "superadmin"
        }
      ]
    }
  ]
}
```

Finally, browse to `/auth` and login with the username and password:

<img src="https://raw.githubusercontent.com/greenpau/caddy-auth-portal/main/assets/docs/images/basic_login.png">

### Password Management

An administrator may change the password directly in
`/etc/caddy/auth/local/users.json` file.

First, download `bcrypt-cli`:

```bash
go get -u github.com/bitnami/bcrypt-cli
```

Then, use it to generate a new password:

```bash
$ echo -n "password123" | bcrypt-cli -c 10
$2a$10$OVnOaHDkcOXfbUZPFh5qt.yJqUt6pl9uJaqEMxxM.vS5fY/cZNmsq
```

Finally, replace the newly generated password is user database file.

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->
## LDAP Authentication Backend

It is recommended reading the documentation for Local backend, because
it outlines important principles of operation of all backends.

Additionally, the LDAP backend works in conjunction with Local backend.
As you will see later, the two can be used together by introducing a
dropdown in UI interface to choose local versus LDAP domain authentication.

The reference configuration for the backend is `assets/conf/ldap/config.json`.

The following Caddy endpoint at `/auth` authentications users
from `contoso.com` domain.

There is a single LDAP server associated with the domain: `ldaps://ldaps.contoso.com`.

The plugin DOES NOT ignore certificate errors when connecting to the servers.
However, one may ignore the errors by setting `ignore_cert_errors` to `true`.

As a better alternative to ignoring certificate errors, the plugin allows
adding trusted certificate authorities via `trusted_authority` Caddyfile directive:

```
          servers {
            ldaps://ldaps.contoso.com
          }
          trusted_authority /etc/gatekeeper/tls/trusted_authority/contoso_com_root1_ca_cert.pem
          trusted_authority /etc/gatekeeper/tls/trusted_authority/contoso_com_root2_ca_cert.pem
          trusted_authority /etc/gatekeeper/tls/trusted_authority/contoso_com_root3_ca_cert.pem
```

The LDAP attribute mapping to JWT fields is as follows.

| **JWT Token Field** | **LDAP Attribute** |
| --- | --- |
| `name` | `givenName` |
| `surname` | `sn` |
| `username` | `sAMAccountName` |
| `member_of` | `memberOf` |
| `email` | `mail` |

The plugin uses `authzsvc` domain user to perform LDAP bind.

The base search DN is `DC=CONTOSO,DC=COM`.

The plugin accepts username (`sAMAccountName`) or email address (`mail`)
and uses the following search filter: `(&(|(sAMAccountName=%s)(mail=%s))(objectclass=user))`.

For example:

```json
      {
        "Name": "sAMAccountName",
        "Values": [
          "jsmith"
        ]
      },
      {
        "Name": "mail",
        "Values": [
          "jsmith@contoso.com"
        ]
      }
```

Upon successful authentication, the plugin assign the following rules
to a user, provided the user is a member of a group:

| **JWT Role** | **LDAP Group Membership** |
| --- | --- |
| `admin` | `CN=Admins,OU=Security,OU=Groups,DC=CONTOSO,DC=COM` |
| `editor` | `CN=Editors,OU=Security,OU=Groups,DC=CONTOSO,DC=COM` |
| `viewer` | `CN=Viewers,OU=Security,OU=Groups,DC=CONTOSO,DC=COM` |

The security of the `password` could be improved by the following techniques:

* pass the password via environment variable `LDAP_USER_SECRET`
* store the password in a file and pass the file inside the `password`
  field with `file:` prefix, e.g. `file:/path/to/password`.

### Configuration Primer

The following `Caddyfile` secures Prometheus/Alertmanager services. Users may access
using local and LDAP credentials.

```
{
  http_port     8080
  https_port    8443
  debug
}

127.0.0.1:8443 {
  route /auth* {
    auth_portal {
      path /auth
      backends {
        local_backend {
          method local
          path assets/conf/local/auth/user_db.json
          realm local
        }
        ldap_backend {
          method ldap
          realm contoso.com
          servers {
            ldaps://ldaps.contoso.com
          }
          trusted_authority /etc/gatekeeper/tls/trusted_authority/contoso_com_root1_ca_cert.pem
          trusted_authority /etc/gatekeeper/tls/trusted_authority/contoso_com_root2_ca_cert.pem
          trusted_authority /etc/gatekeeper/tls/trusted_authority/contoso_com_root3_ca_cert.pem
          attributes {
            name givenName
            surname sn
            username sAMAccountName
            member_of memberOf
            email mail
          }
          username "CN=authzsvc,OU=Service Accounts,OU=Administrative Accounts,DC=CONTOSO,DC=COM"
          # password "P@ssW0rd123"
          password "file:/etc/gatekeeper/auth/ldap.secret"
          search_base_dn "DC=CONTOSO,DC=COM"
          search_filter "(&(|(sAMAccountName=%s)(mail=%s))(objectclass=user))"
          groups {
            "CN=Admins,OU=Security,OU=Groups,DC=CONTOSO,DC=COM" admin
            "CN=Editors,OU=Security,OU=Groups,DC=CONTOSO,DC=COM" editor
            "CN=Viewers,OU=Security,OU=Groups,DC=CONTOSO,DC=COM" viewer
          }
        }
      }
      jwt {
        token_name access_token
        token_secret 0e2fdcf8-6868-41a7-884b-7308795fc286
      }
      ui {
        logo_url "https://caddyserver.com/resources/images/caddy-circle-lock.svg"
        logo_description "Caddy"
        links {
          "Prometheus" /prometheus
          "Alertmanager" /alertmanager
          "My App" /myapp
        }
      }
    }
  }

  route /prometheus* {
    jwt {
      primary yes
      trusted_tokens {
        static_secret {
          token_name access_token
          token_secret 0e2fdcf8-6868-41a7-884b-7308795fc286
        }
      }
      auth_url /auth
      allow roles anonymous guest admin
      allow roles superadmin
      allow roles admin editor viewer
      allow roles AzureAD_Administrator AzureAD_Editor AzureAD_Viewer
    }
    uri strip_prefix /prometheus
    reverse_proxy http://127.0.0.1:9080
  }

  route /alertmanager* {
    jwt
    uri strip_prefix /alertmanager
    reverse_proxy http://127.0.0.1:9083
  }

  route /myapp* {
    jwt
    respond * "myapp" 200
  }

  route /version* {
    respond * "1.0.0" 200
  }

  route {
    redir https://{hostport}/auth 302
  }
}
```

### LDAP Authentication Process

The plugin does not keep connections open to LDAP servers. The plugin
tears a connection down each time it finishes authenticating a request
associated with the connection.

First, the plugin uses `username` and `password` to bind to an LDAP
server. The purpose of the connection is searching for user objects
in the server's directory.

The plugin takes the username provided in a request. Next, the
plugin substitutes `%s` with the username in its search filter, i.e.
`(&(|(sAMAccountName=%s)(mail=%s))(objectclass=user))`.

The plugin initiates a search for a user object in the scope provided
via `search_base_dn`, e.g. `DC=CONTOSO,DC=COM`.

If the number of objects in the result of the search is not `1`, then
authentication fails.

Typically, the response would have the following structure:

```json
[
  {
    "DN": "CN=Smith\\, John,OU=Users,DC=CONTOSO,DC=COM",
    "Attributes": [
      {
        "Name": "sn",
        "Values": [
          "Smith"
        ]
      },
      {
        "Name": "givenName",
        "Values": [
          "John"
        ]
      },
      {
        "Name": "memberOf",
        "Values": [
          "CN=Admins,OU=Security,OU=Groups,DC=CONTOSO,DC=COM",
          "CN=Editors,OU=Security,OU=Groups,DC=CONTOSO,DC=COM",
          "CN=Viewers,OU=Security,OU=Groups,DC=CONTOSO,DC=COM"
        ]
      },
      {
        "Name": "sAMAccountName",
        "Values": [
          "jsmith"
        ]
      },
      {
        "Name": "mail",
        "Values": [
          "jsmith@contoso.com"
        ]
      }
    ]
  }
]
```

The plugin iterates over `memberOf` attribute and compares the
values to its group mapping:

```json
              "groups": [
                {
                  "dn": "CN=Admins,OU=Security,OU=Groups,DC=CONTOSO,DC=COM",
                  "roles": [
                    "admin"
                  ]
                },
                {
                  "dn": "CN=Editors,OU=Security,OU=Groups,DC=CONTOSO,DC=COM",
                  "roles": [
                    "editor"
                  ]
                },
                {
                  "dn": "CN=Viewers,OU=Security,OU=Groups,DC=CONTOSO,DC=COM",
                  "roles": [
                    "viewer"
                  ]
                }
              ]
```

If there are no matches, the authentication fails.

Once the plugin determines the user's roles, e.g. `admin`, `editor`, `viewer`,
the plugin actually checks whether the user's password is valid.

It does so by doing LDAP re-binding with the user's DN and the password provided
in the request. In this example, the user's DN is
`CN=Smith\\, John,OU=Users,DC=CONTOSO,DC=COM`.

If the re-binding is successful, the plugin issues a JWT token.

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->

## SAML Authentication Backend

The plugin supports the following SAML identity providers (IdP):

* Azure Active Directory (Office 365) Applications

If you would like to see the support for the following identity providers,
please reach out:

* Salesforce
* Okta
* Ping Identity

[:arrow_up: Back to Top](#table-of-contents)

### Time Synchronization

Importantly, SAML assertion validation checks timestamps. It is
critical that the application validating the assertions maintains
accurate clock. The out of sync time WILL result in failed
authentications.

### Configuration

The following configuration is common across variations of SAML backend:

```
      backends {
        azure_saml_backend {
          method saml
          realm azure
          provider azure
        }
      }
```

| **Parameter Name** | **Description** |
| --- | --- |
| `method` | Must be set to `saml` |
| `realm` | The realm is used to distinguish between various SAML authentication providers |
| `provider` | It is either `generic` or specific, e.g. `azure`, `okta`, etc. |

The URL for the SAML endpoint is: `<AUTH_PORTAL_PATH>/saml/<REALM_NAME>`.

If you specify `realm` as `azure` and the portal is being served at
`/auth`, then you could access the endpoint via `/auth/saml/azure`.

The Reply URL could be `https://localhost:8443/auth/saml/azure`.

### Azure Active Directory (Office 365) Applications

#### Azure AD SAML Configuration

The Azure SAML backend configuration:

```
      backends {
        azure_saml_backend {
          method saml
          realm azure
          provider azure
          idp_metadata_location /etc/gatekeeper/auth/idp/azure_ad_app_metadata.xml
          idp_sign_cert_location /etc/gatekeeper/auth/idp/azure_ad_app_signing_cert.pem
          tenant_id "1b9e886b-8ff2-4378-b6c8-6771259a5f51"
          application_id "623cae7c-e6b2-43c5-853c-2059c9b2cb58"
          application_name "My Gatekeeper"
          entity_id "urn:caddy:mygatekeeper"
          acs_url https://mygatekeeper/auth/saml/azure
          acs_url https://mygatekeeper.local/auth/saml/azure
          acs_url https://192.168.10.10:3443/auth/saml/azure
          acs_url https://localhost:3443/auth/saml/azure
        }
      }
```

The plugin supports the following parameters for Azure Active
Directory (Office 365) applications:

| **Parameter Name** | **Description** |
| --- | --- |
| `idp_metadata_location` | The url or path to Azure IdP Metadata |
| `idp_sign_cert_location` | The path to Azure IdP Signing Certificate |
| `tenant_id` | Azure Tenant ID |
| `application_id` | Azure Application ID |
| `application_name` | Azure Application Name |
| `entity_id` | Azure Application Identifier (Entity ID) |
| `acs_url` | Assertion Consumer Service URLs |

Use the `acs_url` directive to list all URLs the users of the application
can reach it at. One URL per line:

```
  acs_url https://mygatekeeper/auth/saml/azure
  acs_url https://mygatekeeper.local/auth/saml/azure
  acs_url https://192.168.10.10:3443/auth/saml/azure
  acs_url https://localhost:3443/auth/saml/azure
```

[:arrow_up: Back to Top](#table-of-contents)

#### Set Up Azure AD Application

In Azure AD, you will have an application, e.g. "My Gatekeeper".

The application is a Caddy web server running on port 3443 on
`localhost`. This example meant to emphasize that the authorization
is asynchronious. That is when a user clicks on "My Gatekeeper" icon
in Office 365, the browser takes the user to a sign in page
at URL `https://localhost:3443/saml`.

![Azure AD App Registration - Overview](./assets/docs/images/azure_app_registration_overview.png)

The Application Identifiers are as follows:

* Application (client) ID: `623cae7c-e6b2-43c5-853c-2059c9b2cb58`
* Directory (tenant) ID: `1b9e886b-8ff2-4378-b6c8-6771259a5f51`
* Object ID: `515d2e8b-7548-413f-abee-a23ece1ea576`

The "Branding" page configures "Home Page URL".

![Azure AD App Registration - Branding](./assets/docs/images/azure_app_registration_branding.png)

For demostration purposes, we will create the following "Roles" in the application:

| **Azure Role Name** | **Role Name in SAML Assertion** |
| --- | --- |
| Viewer | AzureAD_Viewer |
| Editor | AzureAD_Editor |
| Administrator | AzureAD_Administrator |

Use "Manifest" tab to add roles in the manifest via `appRoles` key:

![Azure AD App Registration - Manifest - User Roles](./assets/docs/images/azure_app_registration_user_roles.png)

```json
{
  "allowedMemberTypes": [
    "User"
  ],
  "description": "Administrator",
  "displayName": "Administrator",
  "id": "91287df2-7028-4d5f-b5ae-5d489ba217dd",
  "isEnabled": true,
  "lang": null,
  "origin": "Application",
  "value": "AzureAD_Administrator"
},
{
  "allowedMemberTypes": [
    "User"
  ],
  "description": "Editor",
  "displayName": "Editor",
  "id": "d482d827-1757-4f60-9bea-021c10037674",
  "isEnabled": true,
  "lang": null,
  "origin": "Application",
  "value": "AzureAD_Editor"
},
{
  "allowedMemberTypes": [
    "User"
  ],
  "description": "Viewer",
  "displayName": "Viewer",
  "id": "c69f7abd-0a88-401e-b515-92d74b6fff2f",
  "isEnabled": true,
  "lang": null,
  "origin": "Application",
  "value": "AzureAD_Viewer"
}
```

After, we added the roles, we could assign any of the roles to a user:

![Azure AD App - Users and Groups - Add User](./assets/docs/images/azure_app_add_user.png)

The app is now available to the provisioned users in Office 365:

![Office 365 - Access Application](./assets/docs/images/azure_app_user_access.png)

[:arrow_up: Back to Top](#table-of-contents)

#### Configure SAML Authentication

Go to "Enterprise Application" and browse to "My Gatekeeper" application.

There, click "Single Sign-On" and select "SAML" as the authentication method.

![Azure AD App - Enable SAML](./assets/docs/images/azure_app_saml_enable.png)

Next, in the "Set up Single Sign-On with SAML", provide the following
"Basic SAML Configuration":

* Identifier (Entity ID): `urn:caddy:mygatekeeper`
* Reply URL (Assertion Consumer Service URL): `https://localhost:3443/auth/saml/azure`

![Azure AD App - Basic SAML Configuration](./assets/docs/images/azure_app_saml_id.png)

Under "User Attributes & Claims", add the following claims to the list of
default claims:

| **Namespace** | **Claim name** | **Value** |
| --- | --- | --- |
| `http://claims.contoso.com/SAML/Attributes` | `RoleSessionName` | `user.userprincipalname` |
| `http://claims.contoso.com/SAML/Attributes` | `Role` | `user.assignedroles` |
| `http://claims.contoso.com/SAML/Attributes` | `MaxSessionDuration` | `3600` |

![Azure AD App - User Attributes and Claims](./assets/docs/images/azure_app_saml_claims.png)

Next, record the following:
* App Federation Metadata Url
* Login URL

Further, download:
* Federation Metadata XML
* Certificate (Base64 and Raw)

![Azure AD App - SAML Signing Certificate](./assets/docs/images/azure_app_saml_other.png)

[:arrow_up: Back to Top](#table-of-contents)

#### Azure AD IdP Metadata and Certificate

The following command downloads IdP metadata file for Azure AD Tenant with
ID `1b9e886b-8ff2-4378-b6c8-6771259a5f51`. Please note the `xmllint` utility
is a part of `libxml2` library.

```bash
mkdir -p /etc/gatekeeper/auth/saml/idp/
curl -s -L -o /tmp/federationmetadata.xml https://login.microsoftonline.com/1b9e886b-8ff2-4378-b6c8-6771259a5f51/federationmetadata/2007-06/federationmetadata.xml
sudo mkdir -p /etc/gatekeeper/auth/saml/idp/
cat /tmp/federationmetadata.xml | xmllint --format - | sudo tee /etc/gatekeeper/auth/saml/idp/azure_ad_app_metadata.xml
```

The `/etc/gatekeeper/auth/saml/idp/azure_ad_app_metadata.xml` contains IdP metadata.
This file contains the data necessary to verify the SAML claims received by this
service and signed by Azure AD. The `idp_metadata` argument is being used to
pass the location of IdP metadata.

Next, download the "Certificate (Base64)" and store it in
`/etc/gatekeeper/auth/saml/idp/azure_ad_app_signing_cert.pem`.

[:arrow_up: Back to Top](#table-of-contents)

#### User Interface Options

First option is a login button on the login server web page. Once Azure AD has
been enabled, the `/auth/saml/azure` page will have "Sign in with Office 365" button

![Azure AD App - Login with Azure Button](./assets/docs/images/login_with_azure_button.png?width=20px)

Second option is Office 365 applications. When a user click on the
application's icon in Office 365, the user gets redirected to the web
server by Office 365.

![Office 365 - Access Application](./assets/docs/images/azure_app_user_access.png)

The URL is `https://localhost:3443/auth/saml/azure`.

[:arrow_up: Back to Top](#table-of-contents)

#### Development Notes

The below are the headers of the redirected `POST` request that the user's
browser makes upon clicking "My Gatekeeper" application:

```
Method: POST
URL: /auth/saml/azure
Protocol: HTTP/2.0
Host: localhost:3443
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9,ru;q=0.8
Cache-Control: max-age=0
Content-Length: 7561
Content-Type: application/x-www-form-urlencoded
Origin: https://login.microsoftonline.com
Referer: https://login.microsoftonline.com/
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: cross-site
Upgrade-Insecure-Requests: 1
```

The above redirect contains `login.microsoftonline.com` in the request's
`Referer` header. It is the trigger to perform SAML-based authorization.

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->

## OAuth 2.0 and OpenID Connect (OIDC) Authentication Backend

### OAuth 2.0 Flow

Generally, a developer must create an "Application" with and identity provider.
e.g. Google, Okta, Azure, etc., and obtain OAuth 2.0 "Client ID" and
"Client Secret".

The authentication flow begins with obtaining "Authorization Code" from an
identity provider.

* What does the portal send to Okta?
  - Client ID
  - Redirect URI
  - Response Type
  - Scope

* What does the portal receive from Okta?
  - Authorization Code

Once, the portal has the "Authorization Code", it could get "Access Token"
to access the user's data at the identity provider.

* What does the portal send to Okta?
  - Authorization Code
  - Client ID
  - Client Secret

* What does Okta respond with?
  - Access Token

* What could the portal use "Access Token" for?
  -  Make API calls to obtain user information

The OpenID Connect (OIDC) adds login and profile information about the person
who is logged in. The differences between standard OAuth2.0 flow are:

1. In the initial request, a specific scope of `openid` is used
2. In the final exchange the Client receives both an "Access Token" and an "ID Token" (JWT Token).

References:
* [Mozilla - OIDC in a nutshell](https://infosec.mozilla.org/guidelines/iam/openid_connect.html#oidc-in-a-nutshell)

[:arrow_up: Back to Top](#table-of-contents)

### Adding Role Claims

The Caddyfile `user` directive allows adding roles to
a user based on the user's email.

A user with email claim of `contoso.com` would get an additional `superuser` role.

```
          user jsmith@contoso.com add role superuser
```

A user with the email address beginning with `jsmith` would get additional roles.
Specifically, it would be viewer, editor, and admin.

```
          user "^greenpau" regex add roles viewer editor admin
```

All users with `contoso.com` email address would get "employee" role:

```
          user "@contoso.com$" regex add role employee

```

In sum, Caddyfile may look as follows:

```
myapp.localdomain.local, localhost, 127.0.0.1 {
  route /auth* {
    auth_portal {
      path /auth
      backends {
        google_oauth2_backend {
          method oauth2
          realm google
          provider google
          client_id XXXXXXXXXXXXXX.apps.googleusercontent.com
          client_secret YYYYYYYYYYYYYYYYY
          scopes openid email profile
          user "^greenpau" regex add role superuser
        }
```

### OAuth 2.0 Authorization Servers and Identity Providers

The Caddyfile snippet for generic (non-specific) OAuth 2.0 backend.

Based on the below configuration, OAuth 2.0 endpoint is `/auth/oauth2/generic`. If a user
browses to the endpoint, the user will be redirected to the provider discovered via
`metadata_url` and `base_auth_url` URLs.

```
127.0.0.1, localhost {
  route /auth* {
    auth_portal {
      backends {
        generic_oauth2_backend {
          method oauth2
          realm generic
          provider generic
          client_id XXXXXXXXXXXXXXXXXXX
          client_secret YYYYYYYYYYYYYYYYYYYYYY
          scopes openid email profile
        }
```

Next, protect a route, e.g. `/sso/oauth2/generic*`. When a user accesses the page, the
the user will be redirected to `/auth/oauth2/generic` and, then, to the provider. Once authenticated,
the user will be redirected back to `/sso/oauth2/generic...`, i.e. back to the path the user
came from.

```
  route /sso/oauth2/generic* {
    jwt {
      auth_url /auth/oauth2/generic
    }
    respond * "generic oauth2 sso" 200
  }
```

#### Okta

Create an "Application," by browsing to "Applications" and clicking
"Add Application" button.

For a website, the choice is "Web".

![Okta Developer - New App Choice](./assets/docs/images/oauth2_okta_new_app_choice.png)

Provided your application is running on localhost port 8443, provide the following
information:

Base URI entries:
* `https://myapp.localdomain.local:8443/`
* `https://localhost:8443/`

Login redirect URIs:
* `https://myapp.localdomain.local:8443/auth/oauth2/okta/authorization-code-callback`
* `https://localhost:8443/auth/oauth2/okta/authorization-code-callback`

Logout redirect URIs:
* `https://myapp.localdomain.local:8443/auth/logout`
* `https://localhost:8443/auth/logout`

Group Assignments:
* Everyone
* Viewer
* Editor
* Administrator

Grant type allowed:
* Client acting on behalf of itself
  - Client Credentials: No
* Client acting on behalf of a user
  - Authorization Code: Yes
  - Refresh Token: No
  - Implicit (Hybrid): No


![Okta Developer - New App Setup](./assets/docs/images/oauth2_okta_new_app.png)

Review the newly created application.

![Okta Developer - Settings - General](./assets/docs/images/oauth2_okta_app_settings_01.png)

Store the credentials securely.

![Okta Developer - Settings - Client Credentials](./assets/docs/images/oauth2_okta_app_settings_02.png)

Review default Sign On Policy.

![Okta Developer - Settings - Sign On](./assets/docs/images/oauth2_okta_app_settings_03.png)

By default, the default Authorization Server has no `groups` scope.

Therefore, browse to "API", "Authorization Servers" and select "default".

![Okta Developer - API](./assets/docs/images/okta_configure_scope_01.png)

![Okta Developer - API](./assets/docs/images/okta_configure_scope_02.png)

Next, browse to "Scopes" and click "Add Scope".

Fill out the "Add Scope" form:
* Name: `groups`
* Description: `This allows the app to view your group memberships.`
* Check "Set as a default scope"
* Check "Include in public metadata"

![Okta Developer - Add Scope](./assets/docs/images/okta_configure_scope_03.png)

Next, browse to "Claims" and click "Add Claim".

Fill out the "Add Claim" form:
* Name: `groups`
* Include in token type: "ID Token", "Always"
* Value type: Groups
* Filter: Matches `.*` regex
* Include in: The "groups" scope

![Okta Developer - Add Scope](./assets/docs/images/okta_configure_scope_04.png)

Next, review [Okta OpenID Connect and OAuth 2.0 API - Get Started](https://developer.okta.com/docs/reference/api/oidc/#get-started).

The Caddyfile snipper for Okta OAuth 2.0 backend is as follows.

Based on the below configuration, Okta endpoint is `/auth/oauth2/okta`. If a user
browses to the endpoint, the user will be redirected to Okta.

```
127.0.0.1, localhost {
  route /auth* {
    auth_portal {
      backends {
        okta_oauth2_backend {
          method oauth2
          realm okta
          provider okta
          domain_name dev-680653.okta.com
          client_id 0oa121qw81PJW0Tj34x7
          client_secret b3aJC5E59hU18YKC7Yca3994F4qFhWiAo_ZojanF
          server_id default
          scopes openid email profile groups
        }
```

Next, protect a route, e.g. `/sso/oauth2/okta*`. When a user accesses the page, the
the user will be redirected to `/auth/oauth2/okta` and, then, to Okta. Once authenticated,
the user will be redirected back to `/sso/oauth2/okta...`, i.e. back to the path the user
came from.

```
  route /sso/oauth2/okta* {
    jwt {
      auth_url /auth/oauth2/okta
    }
    respond * "okta oauth2 sso" 200
  }
```

Provided the Okta domain is `dev-680653.okta.com`, the authorization server is
`default`, and Client ID is `0oa121qw81PJW0Tj34x7`, check OpenID configuration:

```bash
curl -X GET "https://dev-680653.okta.com/oauth2/default/.well-known/openid-configuration?client_id=0oa121qw81PJW0Tj34x7" | jq
```

By default, the plugin logs public keys from `keys` endpoint.

[:arrow_up: Back to Top](#table-of-contents)

#### Google Identity Platform

References:
* [Google Identity Platform - Using OAuth 2.0 for Web Server Applications](https://developers.google.com/identity/protocols/oauth2/web-server#httprest_2)
* [Google Identity Platform - Identity Platform - OpenID Connect](https://developers.google.com/identity/protocols/oauth2/openid-connect)

The Caddyfile snippet for Google OAuth 2.0 OpenID backend is as follows.

Based on the below configuration, Google endpoint is `/auth/oauth2/google`. If a user
browses to the endpoint, the user will be redirected to Google.

```
127.0.0.1, localhost {
  route /auth* {
    auth_portal {
      backends {
        google_oauth2_backend {
          method oauth2
          realm google
          provider google
          client_id XXXXXXXXXXXXXXXXX.apps.googleusercontent.com
          client_secret YYYYYYYYYYYYYYYYYY
          scopes openid email profile
        }
```

Next, protect a route, e.g. `/sso/oauth2/google*`. When a user accesses the page, the
the user will be redirected to `/auth/oauth2/google` and, then, to Google. Once authenticated,
the user will be redirected back to `/sso/oauth2/google...`, i.e. back to the path the user
came from.

```
  route /sso/oauth2/google* {
    jwt {
      auth_url /auth/oauth2/google
    }
    respond * "google oauth2 sso" 200
  }
```

`[assets/conf/google/Caddyfile](assets/conf/google/Caddyfile)`
contains a full example of setting up Google authentication.

First, create new application, e.g. "My Gatekeeper".

![Google Identity Platform - Identity Platform - New Application](./assets/docs/images/oauth2_google_new_app.png)

After the creation of the app, you will land on Credentials page.

![Google Identity Platform - Identity Platform - Credentials](./assets/docs/images/oauth2_google_credentials.png)

Click "Configure Consent Screen" and select an appropriate option, e.g. "External".

![Google Identity Platform - Identity Platform - Consent Screen](./assets/docs/images/oauth2_google_consent_screen.png)

Next, provide the name for the application, e.g. "My Gatekeeper" and select appropriate support email.

![Google Identity Platform - Identity Platform - Consent Screen Configuration](./assets/docs/images/oauth2_google_consent_screen_config.png)

After configuring the consent screen you will see the following.

![Google Identity Platform - Identity Platform - Consent Screen Verification](./assets/docs/images/oauth2_google_consent_screen_verification.png)

Next, browse to "Credentials" and click "Create Credentials". Then, choose "OAuth client ID":

![Google Identity Platform - Identity Platform - New Credentials](./assets/docs/images/oauth2_google_new_credentials.png)

First, choose the type of credentials:

![Google Identity Platform - Identity Platform - Consent Screen](./assets/docs/images/oauth2_google_new_credentials_type_choice.png)

Next, provide Redirect URL:

![Google Identity Platform - Identity Platform - Consent Screen](./assets/docs/images/oauth2_google_new_credentials_uri_choice.png)

Login redirect URIs:
* `https://localhost:8443/auth/oauth2/google/authorization-code-callback`

Finally, you will get a confirmation. Store the Client ID and Client Secret securely.

![Google Identity Platform - Identity Platform - Consent Screen](./assets/docs/images/oauth2_google_new_credentials_confirm.png)

[:arrow_up: Back to Top](#table-of-contents)

#### LinkedIn

First, browse to https://www.linkedin.com/developers/apps/new and create an application.

![LinkedIn Developers - New Application](./assets/docs/images/oauth2_linkedin_new_app.png)

Next, note the "Client Secret"


![LinkedIn Developers - Auth Screen](./assets/docs/images/oauth2_linkedin_auth_screen.png)


After that, add "redirect URLS":

```
https://localhost:8443/auth/oauth2/linkedin/authorization-code-callback
```

![LinkedIn Developers - Auth Screen - Redirect URLs](./assets/docs/images/oauth2_linkedin_redirect_url.png)

Next, browse to "Products" tab and enabled "Sign In with LinkedIn":

![LinkedIn Developers - Products Screen](./assets/docs/images/oauth2_linkedin_products_screen.png)

References:
* [LinkedIn - LinkedIn API Documentation - Authentication - Authorization Code Flow](https://docs.microsoft.com/en-us/linkedin/shared/authentication/authorization-code-flow)
* [LinkedIn - Consumer Solutions Platform - Integrations - Sign In with LinkedIn](https://docs.microsoft.com/en-us/linkedin/consumer/integrations/self-serve/sign-in-with-linkedin)

#### Auth0

TBD.

[:arrow_up: Back to Top](#table-of-contents)

#### OneLogin

TBD.

[:arrow_up: Back to Top](#table-of-contents)


#### Microsoft

To register an OAuth2 application for login with Microsoft accounts
(either personal, i.e. Live or Xbox accounts, or enterprise,
i.e. AzureAD accounts), you can follow the documentation at
`https://docs.microsoft.com/en-us/advertising/guides/authentication-oauth-identity-platform`.

In summary, you open the `Azure Active Directory` section on
`https://portal.azure.com` and navigate to `App registrations`.

![Azure Active Directory - App registrations](./assets/docs/images/oauth2_azure_new_app.png)

There you select `New registration` and enter your application's name,
your choice of supported account types and the the redirect URI

```
https://localhost:8443/auth/oauth2/azure/authorization-code-callback
```

![Azure Active Directory - App registrations - New application](./assets/docs/images/oauth2_azure_new_application_details.png)

As soon as the application registration was successfully created, you
can note down its `Application (client) ID` listed in the `Essentials`
section at the top.

![Azure Active Directory - App registrations - My Application](./assets/docs/images/oauth2_azure_application_id.png)

Finally, you need to generate a client secret to authenticate
`caddy-auth-portal`. In the sidebar, navigate to `Certificates and
secrets` and click on `New client secret`.

![Azure Active Directory - App registrations - My Application - Certificates and secrets](./assets/docs/images/oauth2_azure_secrets.png)

After the secret was successfully created, copy its value (you won't
be able to retrieve it again!).

![Azure Active Directory - App registrations - My Application - Client secret](./assets/docs/images/oauth2_azure_client_secret.png)

You now have all the necessary information to use the backend in your
`Caddyfile`:

```
        azure_oauth2_backend {
          method oauth2
          realm azure
          provider azure
          client_id 840e455a-69af-47bb-b033-b3a316f1aff0
          client_secret MnE~D8G5Dh_gWRq_jc3uJ8Q8wjBqBv.N3r
          scopes openid email profile
        }
```

If you chose `Accounts in this organizational directory only` as the
account type, you additionally need to add the line

```
          tenant_id <your tenant ID>
```

where the `tenant ID` can either be the actual Directory ID, or its
friendly name `<something>.onmicrosoft.com`.

[:arrow_up: Back to Top](#table-of-contents)

#### Github

Follow the instructions at `https://github.com/settings/apps/new`.

GitHub App name: "My Gatekeeper"

Description: "Caddy v2 Authentication Portal"

Homepage URL: `https://localhost:8443/`

User authorization callback URL: `https://localhost:8443/auth/oauth2/github/authorization-code-callback`

Check "Expire user authorization tokens".

Check "Request user authorization (OAuth) during installation"

Upon successful completion of the instructions, you will get:

![Settings - Developer settings - GitHub Apps - My Gatekeeper](./assets/docs/images/oauth2_github_new_app.png)

Additionally, click "generate a private key" to sign access token requests.

![Settings - Developer settings - GitHub Apps - My Gatekeeper - Private Keys](./assets/docs/images/oauth2_github_sign_keys.png)


Caddyfile configuration:

```
        github_oauth2_backend {
          method oauth2
          realm github
          provider github
          client_id CLIENT_ID
          client_secret CLIENT_SECRET
          scopes user
        }
```

The `github` provider does not have mail claim, i.e. email address. Therefore,
if there is a need to assign a role to a user, one could user `user` directive
to match on `sub`, i.e. subject claim. The `sub` claim is in the format of
`github.com/<GITHUB_HANDLE>`.

```
        github_oauth2_backend {
          ...
          user github.com/greenpau add role superuser
        }
```

The users authenticating via Github would have to accept the terms:

![Settings - Developer settings - GitHub Apps - My Gatekeeper - Accept Terms Screen](./assets/docs/images/oauth2_github_accept_screen.png)

[:arrow_up: Back to Top](#table-of-contents)

#### Facebook

Browse to `https://developers.facebook.com/apps/` and click "Create App".

![Facebook Developers - Apps](./assets/docs/images/oauth2_facebook_apps_screen.png)

When asked about "What do you need your app to do?", choose "Build Connected Experiences".

![Facebook Developers - New App - App Type](./assets/docs/images/oauth2_facebook_apps_type_choice_screen.png)

Next, choose the name for the application:

![Facebook Developers - New App - App Name](./assets/docs/images/oauth2_facebook_apps_name_choice_screen.png)

Once your app (in this case App ID is `38409328409`) is available,
click "Set Up" next to "Facebook Login" product:

![Facebook Developers - App Screen](./assets/docs/images/oauth2_facebook_app_screen.png)

When at Quickstart screen, select "Other".

Next, click "Settings - Advanced" on the left navigation bar and browse to "Security" section.

Set "Require App Secret" to "Yes".
The Client Token is not being used because `client_secret` is being used to calculate `appsecret_proof`.

![Facebook Developers - App Settings - Advanced](./assets/docs/images/oauth2_facebook_app_settings_advanced_screen.png)

Next, click "Settings - Basic" on the left navigation bar and extract "App Secret".
The App Secret is used in `client_secret` Caddyfile directive.

![Facebook Developers - App Settings - Basic](./assets/docs/images/oauth2_facebook_app_settings_basic_screen.png)

Next, click "Settings" under "Facebook Login" on the left navigation bar and browse to "Client OAuth Settings" section:

Set "Valid OAuth Redirect URIs" to:

* `https://localhost:8443/auth/oauth2/facebook/authorization-code-callback`

![Facebook Developers - Facebook Login - Settings](./assets/docs/images/oauth2_facebook_app_login_settings_screen.png)

Additionally, add the URL in "Redirect URI Validator":

* `https://localhost:8443/auth/oauth2/facebook/authorization-code-callback`

The Caddyfile config is as follows:

```
127.0.0.1, localhost {
  route /auth* {
    auth_portal {
      backends {
        facebook_oauth2_backend {
          method oauth2
          realm facebook
          provider facebook
          client_id 38409328409
          client_secret 11899bfcd5745a8ed20235c65638f311
        }
```

When a user get redirected to Facebook Login, the screen looks as follows:

![Facebook Developers - Facebook Login - User Login](./assets/docs/images/oauth2_facebook_user_login_screen.png)

[:arrow_up: Back to Top](#table-of-contents)

<!--- end of section -->

## X.509 Certificate-based Authentication Backend

TBD.

[:arrow_up: Back to Top](#table-of-contents)

## Miscellaneous

### Binding to Privileged Ports

It may be necessary to bind Caddy to privileged port, e.g. 80 or 443.
Grant the `cap_net_bind_service` capability to the Caddy binary, e.g.:

```bash
sudo systemctl stop gatekeeper
sudo rm -rf /usr/local/bin/gatekeeper
sudo cp bin/caddy /usr/local/bin/gatekeeper
sudo setcap cap_net_bind_service=+ep /usr/local/bin/gatekeeper
sudo systemctl start gatekeeper
```

[:arrow_up: Back to Top](#table-of-contents)

### Recording Source IP Address in JWT Token

The `enable source ip tracking` Caddyfile directive instructs
the plugin to record the source IP address when issuing claims.

```
localhost {
  route /auth* {
    auth_portal {
      ...
      enable source ip tracking
      ...
```

This could be useful to force re-authentication when the client IP
address changes.

### Session ID Cache

When the plugin issues JWT tokens, it either passes `jti` values
from upstream providers or generates its own `jti` values.

The plugin stores the mappings between `jti` value and associated
data in a cache. The associated data contains claims and the
metadata from the backend which authenticated a particular session.

This cache is used to assess whether a claim holder is able using
certain portal's capabilities, e.g. add public SSH/GPG key, configure
MFA tokens, change password, etc.

[:arrow_up: Back to Top](#table-of-contents)

### Caddyfile Shortcuts

The following snippet with either `jwt_token_file` or `jwt_token_rsa_file`
Caddyfile directive:

```
    auth_portal {
      jwt_token_file 1 /etc/caddy/auth/jwt/jwt_privatekey.pem
      jwt_token_rsa_file 2 /etc/caddy/auth/jwt/jwt_privatekey.pem
      ...
    }
```

Replaces:

```
    auth_portal {
      jwt {
        token_rsa_file 1 /etc/caddy/auth/jwt/jwt_privatekey.pem
      }
      ...
    }
```

The following snippet with `jwt_token_name` Caddyfile directive:

```
    auth_portal {
      jwt_token_name access_token
      ...
    }
```

Replaces:

```
    auth_portal {
      jwt {
        token_name access_token
      }
      ...
    }
```

The following snippet with `jwt_token_secret` Caddyfile directive:

```
    auth_portal {
      jwt_token_secret bcc8fd6e-8e45-493e-a146-f178ac676841
      ...
    }
```

Replaces:

```
    auth_portal {
      jwt {
        token_secret bcc8fd6e-8e45-493e-a146-f178ac676841
      }
      ...
    }
```

The following snippet with `jwt_token_lifetime` Caddyfile directive:

```
    auth_portal {
      jwt_token_lifetime 3600
      ...
    }
```

Replaces:

```
    auth_portal {
      jwt {
        token_lifetime 3600
      }
      ...
    }
```

[:arrow_up: Back to Top](#table-of-contents)

### Multiple Hosts

`[assets/conf/multihost/Caddyfile](assets/conf/multihost/Caddyfile)`
contains an example of using the same auth config across multiple
hostnames.

[:arrow_up: Back to Top](#table-of-contents)
