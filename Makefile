.PHONY: test ctest covdir coverage docs linter qtest clean dep release logo
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
CADDY_VERSION="v2.1.1"

all:
	@echo "Version: $(PLUGIN_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"
	@mkdir -p bin/
	@rm -rf ./bin/caddy
	@rm -rf ../xcaddy-$(PLUGIN_NAME)/*
	@mkdir -p ../xcaddy-$(PLUGIN_NAME) && cd ../xcaddy-$(PLUGIN_NAME) && \
		xcaddy build $(CADDY_VERSION) --output ../$(PLUGIN_NAME)/bin/caddy \
		--with github.com/greenpau/caddy-auth-portal@$(LATEST_GIT_COMMIT)=$(BUILD_DIR) \
		--with github.com/greenpau/caddy-auth-jwt@latest=$(BUILD_DIR)/../caddy-auth-jwt \
		--with github.com/greenpau/caddy-auth-ui@latest=$(BUILD_DIR)/../caddy-auth-ui \
		--with github.com/greenpau/caddy-request-debug@latest=$(BUILD_DIR)/../caddy-request-debug
	@#bin/caddy run -environ -config assets/conf/local/config.json

linter:
	@echo "Running lint checks"
	@golint -set_exit_status *.go
	@for f in `find ./pkg -type f -name '*.go'`; do echo $$f; go fmt $$f; golint -set_exit_status $$f; done
	@echo "PASS: golint"

test: covdir linter
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out ./*.go

ctest: covdir linter
	@time richgo test $(VERBOSE) $(TEST) -coverprofile=.coverage/coverage.out ./*.go

covdir:
	@echo "Creating .coverage/ directory"
	@mkdir -p .coverage

coverage:
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
	@go test -covermode=count -coverprofile=.coverage/coverage.out ./*.go
	@go tool cover -func=.coverage/coverage.out | grep -v "100.0"

docs:
	@mkdir -p .doc
	@go doc -all > .doc/index.txt
	@#python3 assets/scripts/toc.py > .doc/toc.md

clean:
	@rm -rf .doc
	@rm -rf .coverage
	@rm -rf bin/

qtest:
	@echo "Perform quick tests ..."
	@#go test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestLocalConfig ./*.go
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestLocalCaddyfile ./*.go
	@#go test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestLdapConfig ./*.go
	@#go test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestLdapCaddyfile ./*.go

dep:
	@echo "Making dependencies check ..."
	@go get -u golang.org/x/lint/golint
	@go get -u golang.org/x/tools/cmd/godoc
	@go get -u github.com/kyoh86/richgo
	@go get -u github.com/caddyserver/xcaddy/cmd/xcaddy
	@go get -u github.com/greenpau/versioned/cmd/versioned

release:
	@echo "Making release"
	@go mod tidy
	@go mod verify
	@if [ $(GIT_BRANCH) != "main" ]; then echo "cannot release to non-main branch $(GIT_BRANCH)" && false; fi
	@git diff-index --quiet HEAD -- || ( echo "git directory is dirty, commit changes first" && false )
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
