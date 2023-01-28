KO_DOCKER_REPO ?= kingdonb
TAG ?= latest

tidy:
	go mod tidy -v

build:
	go build ./cmd/kubeconfig-ca-fetch

clean:
	rm -f kubeconfig-ca-fetch

# ko-build:
# 	ko build --local ./cmd/kubeconfig-ca-fetch
#
# ko-publish:
# 	KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build -B ./cmd/gh-app-secret -t $(TAG)
