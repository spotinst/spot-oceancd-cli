# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

VERSION ?= $(shell cat dist/config/VERSION)

COMMIT_SHA ?= $(shell git describe --dirty --always)
DATE_STR ?= $(shell date +%s)
LDFLAGS ?= -X 'spot-oceancd-cli/cmd.commit=$(COMMIT_SHA)' -X 'spot-oceancd-cli/cmd.date=$(DATE_STR)' -X 'spot-oceancd-cli/cmd.version=$(VERSION)'

PREV_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
NEW_HASH=$(shell git rev-parse --verify HEAD)

define newline


endef

define OUTPUT
$(shell git log $(PREV_VERSION)..$(NEW_HASH) --no-merges --pretty=format:'* [view commit](http://github.com/spotinst/spot-oceancd-cli/commit/%H)%s\n' --reverse)

endef

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

test: fmt vet ## Run tests.
	go test ./... -coverprofile cover.out

build: fmt vet test ## Build cli binary.
	go build -ldflags "$(LDFLAGS)" -o dist/oceancd ./

changelog:
	echo "<!-- START ${VERSION} -->" >> "CHANGELOG.md"
	echo "## ${VERSION}" >> "CHANGELOG.md"
	echo "" >> "CHANGELOG.md"
	echo -e '$(subst $(newline),\n,${OUTPUT})' >> "CHANGELOG.md"
	echo "<!-- END ${VERSION} -->" >> "CHANGELOG.md"
	echo "" >> "CHANGELOG.md"