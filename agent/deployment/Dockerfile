FROM openshift/base-centos7

WORKDIR /vpcloue

ADD agent /vpcloue/
ADD config /vpcloue/
RUN chmod +x /vpcloue/agent

EXPOSE 7575
EXPOSE 8080

CMD ["./agent"]