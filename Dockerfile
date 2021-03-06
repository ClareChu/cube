FROM alpine:3.7

EXPOSE 8080
EXPOSE 7575

ENV APP_ROOT=/opt/app-root \
    APP_BIN=${APP_ROOT}/bin \
    PATH=${APP_BIN}:$PATH \
    TZ='Asia/Shanghai'

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} \
     && apk update \
     && apk upgrade \
     && apk --no-cache add ca-certificates iputils\
     && apk add -U tzdata ttf-dejavu busybox-extras curl bash\
     && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
     && adduser -u 1001 -S -G root -g 0 -D -h ${APP_ROOT} -s /sbin/nologin go

# Drop the root user and make the content of /opt/app-root owned by user 1001
RUN chown -R 1001:0 ${APP_ROOT}

WORKDIR ${APP_ROOT}

USER 0