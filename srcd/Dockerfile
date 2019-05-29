FROM golang:1.12-alpine

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH \
    TZ='Asia/Shanghai' \
    GO111MODULE=on \
    DOCKER_API_VERSION=1.24 \
    HOME="/root/app-root" \
    WORKSPACE="/root/workspace"

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} ${HOME} ${WORKSPACE} \
     && apk update \
     && apk upgrade \
     && apk add git \
     && apk add make \
     && apk add bash

WORKDIR  ${HOME}


RUN git clone https://github.com/src-d/enry.git

WORKDIR  ${HOME}/enry

RUN make build
RUN mv ${HOME}/enry/build/bin/enry /usr/local/bin/
RUN    rm -rf ${HOME}/enry

COPY config ${HOME}/config
COPY srcd ${HOME}/srcd

WORKDIR  ${HOME}

EXPOSE 7575
EXPOSE 8080

USER 0

CMD ["./srcd"]
