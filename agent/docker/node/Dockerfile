FROM node:10.15.3-alpine

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH \
    TZ='Asia/Shanghai' \
    GO111MODULE=on \
    DOCKER_API_VERSION=1.24 \
    PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} ${HOME}\
     && apk update \
     && apk upgrade \
     && apk add git \
     && npm install -g yarn@1.16 \
     && yarn config set registry https://registry.npm.taobao.org \
     && npm install -g puppeteer --unsafe-perm=true --allow-root

ADD certs/* /tmp/cert/
RUN apk add --no-cache ca-certificates && \
    mkdir /usr/share/ca-certificates/extra && \
    cp -R /tmp/cert/*.crt /usr/local/share/ca-certificates/ && \
    update-ca-certificates

WORKDIR  ${APP_ROOT}/
COPY config ${APP_ROOT}/config
COPY agent ${APP_ROOT}/agent

EXPOSE 7575
EXPOSE 8080

USER 0

CMD ["./agent"]