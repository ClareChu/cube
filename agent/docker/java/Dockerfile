FROM openjdk:8u212-alpine

EXPOSE 8080

ENV APP_ROOT=/opt/app-root \
    APP_BIN=$APP_ROOT/bin \
    TZ='Asia/Shanghai' \
    MAVEN_HOME='/opt/app-root/apache-maven-3.6.1' \
    PASSWORD=P@55w0dx

RUN  mkdir -p ${APP_BIN} ${APP_ROOT} \
     && apk update \
     && apk upgrade \
     && apk --no-cache add ca-certificates iputils\
     && apk add -U tzdata ttf-dejavu busybox-extras curl bash git curl libc-dev\
     && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
     && adduser -u 1001 -S -G root -g 0 -D -h ${APP_ROOT} -s /sbin/nologin java

WORKDIR ${APP_ROOT}


ADD certs/* /tmp/cert/
RUN apk add --no-cache ca-certificates && \
    mkdir /usr/share/ca-certificates/extra && \
    cp -R /tmp/cert/*.crt /usr/local/share/ca-certificates/ && \
    update-ca-certificates

## install maven

RUN wget http://mirror.bit.edu.cn/apache/maven/maven-3/3.6.1/binaries/apache-maven-3.6.1-bin.tar.gz \
    && tar -xvf apache-maven-3.6.1-bin.tar.gz \
    && curl -u admin:${PASSWORD} https://share.apps.hidevops.io/maven/settings.xml > ${MAVEN_HOME}/conf/settings.xml

USER 0

ENV PATH $MAVEN_HOME/bin:$APP_BIN:$PATH

COPY agent ${APP_ROOT}/agent

COPY config/ ${APP_ROOT}/config/

EXPOSE 7575
EXPOSE 8080

CMD ["./agent"]