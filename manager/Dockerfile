FROM alpine:3.7

EXPOSE 8080
EXPOSE 7575

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH \
    TZ='Asia/Shanghai'

ADD certs /tmp/cert/

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} \
     && apk update \
     && apk upgrade \
     && apk --no-cache add ca-certificates iputils\
     && apk add -U tzdata ttf-dejavu busybox-extras curl bash\
     && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
     && adduser -u 1001 -S -G root -g 0 -D -h ${APP_ROOT} -s /sbin/nologin go

RUN mkdir /usr/share/ca-certificates/extra && \
    cp -R /tmp/cert/*.crt /usr/local/share/ca-certificates/ && \
    update-ca-certificates

COPY ./manager ${APP_BIN}
COPY ./config ${APP_ROOT}/config
COPY ./run.sh ${APP_ROOT}
# Drop the root user and make the content of /opt/app-root owned by user 1001
RUN chown -R 1001:0 ${APP_ROOT}

RUN chmod 755 ${APP_BIN}/manager

WORKDIR ${APP_ROOT}

USER 0

ENTRYPOINT /bin/sh -c ./run.sh