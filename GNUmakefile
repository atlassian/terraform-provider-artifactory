default: build

build: fmtcheck
	go install

test:
	@echo "==> Starting unit tests"
	go test ./... -timeout=30s -parallel=4

artifactory:
	@echo "==> Launching Artifactory in Docker..."
	@scripts/run-artifactory.sh

docker:
	@docker build -t dillongiacoppo/terraform-artifactory .

testacc: fmtcheck artifactory
	@echo "==> Starting integration tests"
	TF_ACC=1 ARTIFACTORY_USERNAME=admin ARTIFACTORY_PASSWORD=password ARTIFACTORY_URL=http://localhost:8080/artifactory \
	go test ./... -v -parallel 20 $(TESTARGS) -timeout 120m
	@docker stop artifactory

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w pkg/artifactory
	goimports -w pkg/artifactory

fmtcheck:
	@echo "==> Checking that code complies with gofmt requirements..."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: build test testacc fmt