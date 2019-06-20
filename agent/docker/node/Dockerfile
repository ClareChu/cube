FROM node:11-alpine

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH \
    TZ='Asia/Shanghai' \
    GO111MODULE=on \
    DOCKER_API_VERSION=1.24

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} ${HOME}\
     && apk update \
     && apk upgrade \
     && apk add git  \
     && apk add --no-cache python2.7

WORKDIR  ${APP_ROOT}/
COPY config ${APP_ROOT}/config
COPY agent ${APP_ROOT}/agent

EXPOSE 7575
EXPOSE 8080

USER 0

CMD ["./agent"]