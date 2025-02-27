# 使用一个基础镜像

FROM golang:1.23 AS builder

ARG GO_LDFLAGS

WORKDIR /go/src/kubesphere.io/edgewize-ota

ADD . /go/src/kubesphere.io/edgewize-ota/

RUN go env -w GOPROXY=https://goproxy.cn

RUN go mod tidy

RUN go build -v -o /usr/local/bin/edgewize-ota-server -ldflags "$GO_LDFLAGS -w -s" edgewize-ota/cmd/edgewize-ota-server/edgewize-ota-server.go


FROM alpine/ansible:2.17.0
RUN apk add --no-cache ca-certificates wget
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub
ENV GLIBC_VERSION=2.34-r0
RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-${GLIBC_VERSION}.apk && \
    apk add --no-cache glibc-${GLIBC_VERSION}.apk --force-overwrite&& \
    rm glibc-${GLIBC_VERSION}.apk


RUN mkdir -p /ansible/roles/edgewize-ota-agent/files/x86_64 /ansible/roles/edgewize-ota-agent/files/arm64 && \
    wget https://github.com/oras-project/oras/releases/download/v1.2.2/oras_1.2.2_linux_amd64.tar.gz && \
    tar -xzf oras_1.2.2_linux_amd64.tar.gz -C /ansible/roles/edgewize-ota-agent/files/x86_64/ oras && \
    rm oras_1.2.2_linux_amd64.tar.gz && \
    chmod +x /ansible/roles/edgewize-ota-agent/files/x86_64/oras && \
    wget https://github.com/oras-project/oras/releases/download/v1.2.2/oras_1.2.2_linux_arm64.tar.gz && \
    tar -xzf oras_1.2.2_linux_arm64.tar.gz -C /ansible/roles/edgewize-ota-agent/files/arm64/ oras && \
    rm oras_1.2.2_linux_arm64.tar.gz && \
    chmod +x /ansible/roles/edgewize-ota-agent/files/arm64/oras

RUN apk add --no-cache sshpass
RUN apk add tini
WORKDIR /ansible
ADD ./edgewize-ota/swagger-ui /ansible/swagger-ui
ADD ./edgewize-ota/roles /ansible/roles
COPY --from=builder /usr/local/bin/edgewize-ota-server /usr/local/bin/edgewize-ota-server
ENV ANSIBLE_HOST_KEY_CHECKING=False

EXPOSE 8082
