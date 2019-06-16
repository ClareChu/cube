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
     && apk add --no-cache python2.7 \
     && python2.7 -m ensurepip && \
        rm -r /usr/lib/python*/ensurepip && \
        pip install --upgrade pip setuptools && \
            if [ ! -e /usr/bin/pip ]; then ln -s pip /usr/bin/pip ; fi && \
            if [[ ! -e /usr/bin/python ]]; then ln -sf /usr/bin/python2.7 /usr/bin/python; fi && \
            rm -r /root/.cache

WORKDIR  ${APP_ROOT}/
COPY config ${APP_ROOT}/config
COPY agent ${APP_ROOT}/agent

EXPOSE 7575
EXPOSE 8080

USER 0

CMD ["./agent"]