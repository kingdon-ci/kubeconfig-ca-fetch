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

tldr: kube.config
	chmod 600 kube.config
	export KUBECONFIG=`pwd`/kube.config
	kubelogin
	kubectl get nodes

supertldr: kube.config
	# !! Overwrites your ~/.kube/config (Ctrl+C to abort)
	sleep 3
	chmod 600 kube.config
	mv kube.config $(HOME)/.kube/config
	kubelogin
	kubectl get nodes
