FROM debian:stretch-slim

# https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl
RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends apt-transport-https ca-certificates dirmngr gnupg; \
	export GNUPGHOME="$(mktemp -d)"; \
	gpg --batch --keyserver ha.pool.sks-keyservers.net --recv-keys '54A6 47F9 048D 5688 D7DA  2ABE 6A03 0B21 BA07 F4FB'; \
	gpg --batch --export --armor '54A6 47F9 048D 5688 D7DA  2ABE 6A03 0B21 BA07 F4FB' > /etc/apt/trusted.gpg.d/kubernetes.gpg.asc; \
	echo 'deb https://apt.kubernetes.io kubernetes-stretch main' > /etc/apt/sources.list.d/kubernetes.list; \
	apt-get update; \
	apt-get install -y --no-install-recommends kubectl=1.6.0*; \
	rm -rf /var/lib/apt/lists/*; \
	kubectl --help

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR $GOPATH/src/kubecon-demo

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		golang-go gcc libc6-dev \
	; \
	rm -rf /var/lib/apt/lists/*

COPY *.go ./
RUN go install -v ./...

COPY static static

CMD ["kubecon-demo"]
