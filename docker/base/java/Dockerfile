FROM openjdk:8u212-alpine

EXPOSE 8080

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
     && adduser -u 1001 -S -G root -g 0 -D -h ${APP_ROOT} -s /sbin/nologin java



WORKDIR ${APP_ROOT}

USER 0