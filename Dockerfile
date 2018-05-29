FROM tianon/docker-tianon

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		golang-go gcc libc6-dev \
	; \
	rm -rf /var/lib/apt/lists/*

WORKDIR $GOPATH/src/dockercon-demo

COPY *.go ./
RUN go install -v ./...

COPY static static

CMD ["dockercon-demo"]
