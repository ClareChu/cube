FROM golang:1.12-alpine

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH \
    TZ='Asia/Shanghai' \
    GO111MODULE=on \
    DOCKER_API_VERSION=1.24

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} ${HOME}\
     && apk update \
     && apk upgrade \
     && apk add git

WORKDIR  ${HOME}/
COPY config ${HOME}/config
COPY agent ${HOME}/agent

EXPOSE 7575
EXPOSE 8080

USER 0

CMD ["./agent"]