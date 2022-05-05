AWSCLI				:= aws
SAMCLI				:= sam

PACKAGE				:= github.com/willvk/go-demo
APP_NAME			:= meetups

GOLANGCI_VERSION	:= 1.43.0
GO_DOCKER_VERSION	:= 1.16
BINDIR				?= $(shell pwd)/bin
GIT_HASH			?= $(shell git rev-parse --short HEAD)

BUILD_OVERRIDES = \
	-X "$(PACKAGE)/internal/app.Name=$(APP_NAME)" \
	-X "$(PACKAGE)/internal/app.BuildDate=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')" \
	-X "$(PACKAGE)/internal/app.Commit=$(GIT_HASH)" \
# the -w -s flags make the binary a bit smaller and
# trimpath shortens build paths in stack traces
LDFLAGS := -ldflags='-w -s $(BUILD_OVERRIDES)' -trimpath
# https://tip.golang.org/cmd/go/#hdr-Module_configuration_for_non_public_modules
# for all modules in the jumacloud org avoid the GOPROXY/GOSUM requests and go direct

S3_BUCKET := $(shell $(AWSCLI) ssm get-parameter \
	--name '/config/ArtifactBucketName' \
	--with-decryption --query 'Parameter.Value' --output text 2> /dev/null || echo NO_AWS_AUTH)

####################
# Package commands #
####################
clean:
	rm -rf $(CURDIR)/dist
.PHONY: clean

lint: $(BINDIR)/golangci-lint generate
	@echo "--- lint all the things"
	@$(BINDIR)/golangci-lint run
.PHONY: lint

validate:
	@echo "--- validate all the things..."
	@cfn-lint --ignore-checks W2001 -- sam/api.sam.yaml
.PHONY: validate

bundle:
	@echo "--- package binary into zipped handler..."
	@cd $(CURDIR)/dist && zip -r -q ../handler.zip .
.PHONY: bundle

###################
#  Test commands  #
###################
test-swagger:
	docker run --rm -it -v $(CURDIR):/var/task stoplight/spectral:4.2 lint -r /var/task/.spectral.yaml -s api-servers -F warn --verbose /var/task/openapi/meetup.yaml
.PHONY: test-swagger

test: generate
	@echo "--- test all the things"
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.txt -covermode count ./...
	@go tool cover -func coverage/coverage.txt | grep total | awk '{print $3}'
.PHONY: test

##################
#  (CI) #
##################
git-config:
	git config --global url.git@github.com:willvk.insteadOf https://github.com/willvk
	mkdir -p ~/.ssh/
	ssh-keyscan -t rsa,dsa -H github.com >> ~/.ssh/known_hosts

docker-ci:
	@echo "docker run..."
	@docker run --rm \
		-v $(shell pwd):/src \
		-v $(SSH_AUTH_SOCK):/ssh-agent:ro -e SSH_AUTH_SOCK=/ssh-agent \
		-w /src -t golang:$(GO_DOCKER_VERSION) make clean git-config build \
		-w $(shell ./dist/meetup)
.PHONY: docker-ci

###################
# Build commands  #
###################
$(BINDIR)/golangci-lint: $(BINDIR)/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint

$(BINDIR)/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b ./bin v${GOLANGCI_VERSION}
	@mv $(BINDIR)/golangci-lint $@

generate:
	@echo "--- generate all the things"
	@go generate ./...
.PHONY: generate

build: generate
	@echo "--- build all the things"
#	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a $(LDFLAGS) -o dist/ ./cmd/meetup/
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a $(LDFLAGS) -o dist/ ./cmd/meetup/

###################
# Deploy commands #
###################
deploy:

