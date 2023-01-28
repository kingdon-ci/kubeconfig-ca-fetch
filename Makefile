KO_DOCKER_REPO ?= kingdonb
TAG ?= latest

all: tidy build kube.config

tidy:
	go mod tidy -v

build:
	go build ./cmd/kubeconfig-ca-fetch

kube.config:
	./kubeconfig-ca-fetch > kube.config

clean:
	rm -f kubeconfig-ca-fetch
	rm -f kube.config

# ko-build:
# 	ko build --local ./cmd/kubeconfig-ca-fetch
#
# ko-publish:
# 	KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build -B ./cmd/gh-app-secret -t $(TAG)
