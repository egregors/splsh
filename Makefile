PROJECT_NAME := "splsh"
PKG := "github.com/egregors/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

CC := gcc 
CGO_ENABLED := 1 
GOOS := darwin
GOARCH := amd64


.PHONY: all build clean test lint docker run

all: run

lint:  ## Lint the files
	@golangci-lint run --config .golangci.yml ./...

test:  ## Run unittests
	@go test -short ${PKG_LIST}

race:  ## Run data race detector
	@go test -race -short ${PKG_LIST}

build:  ## Build the binary file
	# @go build -v $(PKG)
	@go build -gccgoflags="-DDARWIN -x objective-c -fobjc-arc" -ldflags="framework=Cocoa" -v $(PKG)

release: clean build  ## Make MacOS bundle
	@mkdir -p dist
	@mkdir -p dist/Splsh.app
	@mkdir dist/Splsh.app/Contents
	@mkdir -p dist/Splsh.app/Contents/MacOS
	@mkdir -p dist/Splsh.app/Contents/Resources
	@cp ./Splsh dist/Splsh.app/Contents/MacOS
	@cp ./Info.plist dist/Splsh.app/Contents/Info.plist
	@cp ./icon/icon.icns dist/Splsh.app/Contents/Resources/Splsh.icns

sync:  ## Sync deps
	@go mod tidy
	@go mod vendor

run:  ## Go run in debug mode with a race detector
	@go run -race main.go --dbg

clean:  ## Remove previous build
	@rm -f $(PROJECT_NAME)
	@rm -rf dist

docker: ## build Docker image
	@docker build -t splsh .

## Help

help:  ## Show help message
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
	printf "%s\n\n" "Usage: make [task]"; \
	printf "%-20s %s\n" "task" "help" ; \
	printf "%-20s %s\n" "------" "----" ; \
	for help_line in $${help_lines[@]}; do \
		IFS=$$':' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf '\033[36m'; \
		printf "%-20s %s" $$help_command ; \
		printf '\033[0m'; \
		printf "%s\n" $$help_info; \
	done
