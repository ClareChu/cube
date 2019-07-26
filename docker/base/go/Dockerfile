FROM golang:1.12-alpine

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH \
    TZ='Asia/Shanghai' \
    GO111MODULE=on \
    DOCKER_API_VERSION=1.24 \


RUN  mkdir -p ${APP_BIN} ${APP_ROOT} ${HOME}\
     && apk update \
     && apk upgrade \
     && apk add git \
     && apk add curl

RUN apk add --no-cache ca-certificates \
    && curl -u admin:P@55w0dx https://share.apps.hidevops.io/images/dlv?download=true -o dlv \
    && mv dlv /usr/local/bin/

WORKDIR  ${HOME}/

COPY debug /usr/local/bin/

RUN chmod 755 /usr/local/bin/debug


EXPOSE 7575
EXPOSE 8080
