FROM openjdk:8-jre-alpine3.9

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

RUN wget http://mirror.bit.edu.cn/apache/tomcat/tomcat-8/v8.5.41/bin/apache-tomcat-8.5.41.tar.gz \
    && tar -xvf apache-tomcat-8.5.41.tar.gz \
    && mv apache-tomcat-8.5.41 tomcat \
    && rm -rf apache-tomcat-8.5.41.tar.gz

USER 0

CMD ["/opt/app-root/tomcat/bin/catalina.sh", "run"]