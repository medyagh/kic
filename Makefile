GOLINT_VERSION ?= v1.17.1
GOLINT_OPTIONS = --deadline 4m \
	  --build-tags "${MINIKUBE_INTEGRATION_BUILD_TAGS}" \
	  --enable goimports,gocritic,golint,gocyclo,interfacer,misspell,nakedret,stylecheck,unconvert,unparam \
	  --exclude 'variable on range scope.*in function literal|ifElseChain' 


out/linters/golangci-lint:
	mkdir -p out/linters
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b out/linters $(GOLINT_VERSION)


lint: out/linters/golangci-lint
	./out/linters/golangci-lint run ${GOLINT_OPTIONS} ./...


out/e2e:
	GO111MODULE=on go mod download
	mkdir -p out
	GO111MODULE=on go build -o ./out/e2e example/single_node/main.go 
