.PHONY: test ctest covdir coverage docs linter qtest clean dep release logo templates info license
PLUGIN_NAME="caddy-auth-portal"
PLUGIN_VERSION:=$(shell cat VERSION | head -1)
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
LATEST_GIT_COMMIT:=$(shell git log --format="%H" -n 1 | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
BUILD_DIR:=$(shell pwd)
VERBOSE:=-v
ifdef TEST
	TEST:="-run ${TEST}"
endif
CADDY_VERSION="v2.3.0"

all: build

build: info templates license
	@mkdir -p bin/
	@rm -rf ./bin/caddy
	@rm -rf ../xcaddy-$(PLUGIN_NAME)/*
	@mkdir -p ../xcaddy-$(PLUGIN_NAME) && cd ../xcaddy-$(PLUGIN_NAME) && \
		xcaddy build $(CADDY_VERSION) --output ../$(PLUGIN_NAME)/bin/caddy \
		--with github.com/greenpau/caddy-auth-portal@$(LATEST_GIT_COMMIT)=$(BUILD_DIR) \
		--with github.com/greenpau/caddy-auth-jwt@latest=$(BUILD_DIR)/../caddy-auth-jwt \
		--with github.com/greenpau/go-identity@latest=$(BUILD_DIR)/../go-identity \
		--with github.com/greenpau/caddy-trace@latest=$(BUILD_DIR)/../caddy-trace
	@#bin/caddy run -environ -config assets/conf/local/config.json
	@echo "build: OK"

info:
	@echo "Version: $(PLUGIN_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"

templates:
	@./assets/scripts/generate_ui.sh

linter:
	@echo "Running lint checks"
	@golint -set_exit_status *.go
	@for f in `find ./pkg -type f -name '*.go'`; do echo $$f; go fmt $$f; golint -set_exit_status $$f; done
	@echo "PASS: golint"

test: templates license covdir linter docs
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out ./...

ctest: templates license covdir linter
	@time richgo test $(VERBOSE) $(TEST) -coverprofile=.coverage/coverage.out ./...

covdir:
	@echo "Creating .coverage/ directory"
	@mkdir -p .coverage

coverage:
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
	@go test -covermode=count -coverprofile=.coverage/coverage.out ./...
	@go tool cover -func=.coverage/coverage.out | grep -v "100.0"

docs:
	@mkdir -p .doc
	@cat ./assets/docs/pages/*.md > README.md
	@versioned -toc
	@#go doc -all > .doc/index.txt
	@#python3 assets/scripts/toc.py > .doc/toc.md

clean:
	@rm -rf .doc
	@rm -rf .coverage
	@rm -rf bin/

qtest:
	@echo "Perform quick tests ..."
	@#time richgo test  $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestLocalConfig ./*.go
	@#time richgo test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestLocalCaddyfile ./*.go
	@#go test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestLdapConfig ./*.go
	@#go test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestLdapCaddyfile ./*.go
	@#go test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestSamlCaddyfile ./*.go
	@#time richgo test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestGetSourceAddress ./*.go
	@#time richgo test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestNewUserInterface ./pkg/ui/*.go
	@#time richgo test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestCookies ./pkg/cookies/*.go
	@#time richgo test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestCookieLifetime ./*.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestSandboxCache ./pkg/cache/sandbox*.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewSandboxCache ./pkg/cache/sandbox*.go
	@#time richgo test -v -coverprofile=.coverage/coverage.out -run TestNewSandboxHurdle ./pkg/cache/sandbox*.go
	@time richgo test -v -coverprofile=.coverage/coverage.out ./pkg/cache/sandbox*.go
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html

dep:
	@echo "Making dependencies check ..."
	@go get -u golang.org/x/lint/golint
	@go get -u golang.org/x/tools/cmd/godoc
	@go get -u github.com/kyoh86/richgo
	@go get -u github.com/caddyserver/xcaddy/cmd/xcaddy
	@go get -u github.com/greenpau/versioned/cmd/versioned
	@go get -u github.com/google/addlicense

license: dep
	@addlicense -c "Paul Greenberg greenpau@outlook.com" -y 2020 pkg/*/*/*.go pkg/*/*.go *.go

mod:
	@echo "Verifying modules"
	@go mod tidy
	@go mod verify

release: info mod license
	@echo "Making release"
	@if [ $(GIT_BRANCH) != "main" ]; then echo "cannot release to non-main branch $(GIT_BRANCH)" && false; fi
	@git diff-index --quiet HEAD -- || ( echo "git directory is dirty, commit changes first" && git status && false )
	@versioned -patch
	@echo "Patched version"
	@git add VERSION
	@git commit -m "released v`cat VERSION | head -1`"
	@git tag -a v`cat VERSION | head -1` -m "v`cat VERSION | head -1`"
	@git push
	@git push --tags
	@@echo "If necessary, run the following commands:"
	@echo "  git push --delete origin v$(PLUGIN_VERSION)"
	@echo "  git tag --delete v$(PLUGIN_VERSION)"

logo:
	@convert -background black -fill white -font DejaVu-Sans-Bold -size 640x320! -gravity center -pointsize 96 label:'caddy.auth\nportal' PNG32:assets/docs/images/logo.png
