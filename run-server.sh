#!/usr/bin/env bash
set -Eeuo pipefail

[ -d "$HOME/.kube" ]
docker run -dit \
	--name kubecon-demo \
	--network host \
	--restart always \
	-v "$HOME/.kube":/root/.kube:ro \
	172.23.0.151:5000/dockercon-demo
