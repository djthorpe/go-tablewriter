# Paths to tools needed in dependencies
GO := $(shell which go)

# Build flags
BUILD_MODULE := $(shell go list)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitSource=${BUILD_MODULE}
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitTag=$(shell git describe --tags)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitBranch=$(shell git name-rev HEAD --name-only --always)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitHash=$(shell git rev-parse HEAD)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BUILD_FLAGS = -ldflags "-s -w $(BUILD_LD_FLAGS)" 

# Targets
all: dependencies clean test

test:
	@${GO} mod tidy
	@${GO} test -v .

FORCE:

dependencies:
	@test -f "${GO}" && test -x "${GO}"  || (echo "Missing go binary" && exit 1)

clean:
	@echo Clean
	@${GO} mod tidy
	@${GO} clean
