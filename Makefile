.PHONY: help
.PHONY: bump-major bump-minor bump-patch show-version
.PHONY: build push
.PHONY: docker docker-major docker-minor docker-patch

VERSION := $(shell cat VERSION)
DOCKER_REPO := tinywaves/gohub

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  help          - Display this help"
	@echo "  bump-major    - Bump the major version (e.g., 1.2.3 -> 2.0.0)"
	@echo "  bump-minor    - Bump the minor version (e.g., 1.2.3 -> 1.3.0)"
	@echo "  bump-patch    - Bump the patch version (e.g., 1.2.3 -> 1.2.4)"
	@echo "  show-version  - Show current version"
	@echo "  build         - Build Docker image"
	@echo "  push          - Push Docker image to Docker Hub"
	@echo "  docker        - Build and push Docker image to Docker Hub"
	@echo "  docker-major  - Bump major version, build and push Docker image"
	@echo "  docker-minor  - Bump minor version, build and push Docker image"
	@echo "  docker-patch  - Bump patch version, build and push Docker image"

bump-major:
	@$(eval NEW_VERSION := $(shell echo $(VERSION) | awk -F. 'BEGIN{OFS="."} {$$1=$$1+1; $$2=0; $$3=0; print $$0}'))
	@echo $(NEW_VERSION) > VERSION
	@echo "Version bumped to $(NEW_VERSION)"

bump-minor:
	@$(eval NEW_VERSION := $(shell echo $(VERSION) | awk -F. 'BEGIN{OFS="."} {$$2=$$2+1; $$3=0; print $$0}'))
	@echo $(NEW_VERSION) > VERSION
	@echo "Version bumped to $(NEW_VERSION)"

bump-patch:
	@$(eval NEW_VERSION := $(shell echo $(VERSION) | awk -F. 'BEGIN{OFS="."} {$$3=$$3+1; print $$0}'))
	@echo $(NEW_VERSION) > VERSION
	@echo "Version bumped to $(NEW_VERSION)"

show-version:
	@echo "Current version: $(VERSION)"

build:
	@echo "Building Docker image..."
	@rm -f gohub || true
	@GOOS=linux GOARCH=arm go build -o gohub .
	@docker build -t $(DOCKER_REPO):$(shell cat VERSION) .
	@docker tag $(DOCKER_REPO):$(shell cat VERSION) $(DOCKER_REPO):latest
	@rm -f gohub || true

push:
	@echo "Pushing Docker image to Docker Hub..."
	@docker push $(DOCKER_REPO):$(shell cat VERSION)
	@docker push $(DOCKER_REPO):latest
	@docker rmi -f $(DOCKER_REPO):$(shell cat VERSION) $(DOCKER_REPO):latest || true

docker: build push

docker-major: bump-major docker

docker-minor: bump-minor docker

docker-patch: bump-patch docker