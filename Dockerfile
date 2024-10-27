FROM alpine:latest AS binaries
RUN apk --no-cache add wget
RUN mkdir /tmp/binaries
RUN wget -qO /tmp/binaries/bosh "https://github.com/cloudfoundry/bosh-cli/releases/download/v7.8.1/bosh-cli-7.8.1-linux-amd64" && \
  chmod +x /tmp/binaries/bosh

FROM golang:1.23 AS resource
WORKDIR /go/src/github.com/taylorsilva/bosh-release-resource
COPY --from=binaries /tmp/binaries /usr/local/bin
COPY . .
ENV CGO_ENABLED=0
RUN mkdir -p /opt/resource

# RUN git config --global user.email root@localhost
# RUN git config --global user.name root
# RUN go test ./...

RUN git rev-parse HEAD | tee /opt/resource/version
RUN go build -o /opt/resource/check ./check
RUN go build -o /opt/resource/in ./in
RUN go build -o /opt/resource/out ./out

FROM alpine:latest
RUN apk --no-cache add bash ca-certificates curl git openssh-client
COPY --from=binaries /tmp/binaries /usr/local/bin
COPY --from=resource /opt/resource /opt/resource
ADD tasks/create-dev-release tasks/load-release-notes /usr/local/bin/
