# Mio - CI/CD Platform

<p align="center">
  <a href="https://travis-ci.org/hidevopsio/mio?branch=master">
    <img src="https://travis-ci.org/hidevopsio/mio.svg?branch=master" alt="Build Status"/>
  </a>
  <a class="badge-align" href="https://www.codacy.com/app/john-deng/mio?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=hidevopsio/mio&amp;utm_campaign=Badge_Grade"><img src="https://api.codacy.com/project/badge/Grade/ee8ddbf56ece4f46a6efeb216c351a0f"/></a>
  <a href="https://github.com/hidevopsio/mio">
    <img src="https://tokei.rs/b1/github/hidevopsio/mio" />
  </a>
  <a href="https://codecov.io/gh/hidevopsio/mio">
    <img src="https://codecov.io/gh/hidevopsio/mio/branch/master/graph/badge.svg" />
  </a>
  <a href="https://opensource.org/licenses/Apache-2.0">
      <img src="https://img.shields.io/badge/License-Apache%202.0-green.svg" />
  </a>
  <a href="https://goreportcard.com/report/hidevops.io/mio">
      <img src="https://goreportcard.com/badge/hidevops.io/mio" />
  </a>
  <a href="https://godoc.org/hidevops.io/mio">
      <img src="https://godoc.org/github.com/golang/gddo?status.svg" />
  </a>
  <a href="https://gitter.im/hidevopsio/mio">
      <img src="https://img.shields.io/badge/GITTER-join%20chat-green.svg" />
  </a>
</p>

## About

Mio is a Continuous Integration and Continuous Delivery Platform built on container technology.

## 安装

唯一要求需要[GO环境](https://golang.org/)

```bash
go get -u github.com/hidevopsio/mio
```

如果使用GO版本在 `go1.11` 一下 请自行安装dep包管理工具。
在项目下执行

```bash
dep ensure -v
```

如果使用1.11 以上版本最则使用go modules
该项目下分为`console` 和 `node` 俩个项目, console 主要负责cicd 调度 启动 任务分发 等作用， node 主要负责 代码编译 代码测试 镜像的制作等主要工作。

### console镜像制作

```bash
cd console
### 执行build.sh
./build.sh
###脚本大致内容

## console 代码打包linux 
 GOOS=linux go build -o hiadmin

## 镜像制作

docker build -t docker-registry-default.app.vpclub.io/demo/hiadmin:v1 .

docker push docker-registry-default.app.vpclub.io/demo/hiadmin:v1

```

如果使用的是`openshift` 我们提供了 deployment.yml

```yml
apiVersion: apps.openshift.io/v1
kind: DeploymentConfig
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: >
      {"apiVersion":"apps.openshift.io/v1","kind":"DeploymentConfig","metadata":{"annotations":{},"labels":{"app":"hiadmin"},"name":"hiadmin","namespace":"demo"},"spec":{"replicas":1,"selector":{"app":"hiadmin","deploymentconfig":"hiadmin"},"strategy":{"activeDeadlineSeconds":21600},"template":{"spec":{"containers":[{"env":[{"name":"APP_PROFILES_ACTIVE","value":"dev"},{"name":"SCM_URL","value":"http://gitlab.vpclub.cn:8022"}],"image":"docker-registry.default.svc:5000/demo/hiadmin@sha256:6cec6fbb1afed87f60d04050c08248dd3e386302e4e3757931d1625514b176c5","imagePullPolicy":"Always","name":"hiadmin","ports":[{"containerPort":7575,"protocol":"TCP"},{"containerPort":8080,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","volumeMounts":[{"mountPath":"/var/lib/docker","name":"volume1"},{"mountPath":"/var/run/docker.sock","name":"volume2"}]}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30,"volumes":[{"hostPath":{"path":"/var/lib/docker","type":""},"name":"volume1"},{"hostPath":{"path":"/var/run/docker.sock","type":""},"name":"volume2"}]}},"test":false,"triggers":[{"type":"ConfigChange"},{"imageChangeParams":{"automatic":true,"containerNames":["hiadmin"],"from":{"kind":"ImageStreamTag","name":"hiadmin:v1","namespace":"demo"},"lastTriggeredImage":"docker.vpclub.cn/demo/admin:v1"},"type":"ImageChange"}]}}
    openshift.io/generated-by: OpenShiftWebConsole
  creationTimestamp: '2018-10-21T07:06:25Z'
  generation: 188
  labels:
    app: hiadmin
  name: hiadmin
  namespace: demo
  resourceVersion: '139885254'
  selfLink: /apis/apps.openshift.io/v1/namespaces/demo/deploymentconfigs/hiadmin
  uid: d2bcdad9-d4ff-11e8-bb8c-005056935c80
spec:
  replicas: 1
  selector:
    app: hiadmin
    deploymentconfig: hiadmin
  strategy:
    activeDeadlineSeconds: 21600
    resources: {}
    rollingParams:
      intervalSeconds: 1
      maxSurge: 25%
      maxUnavailable: 25%
      timeoutSeconds: 600
      updatePeriodSeconds: 1
    type: Rolling
  template:
    metadata:
      annotations:
        openshift.io/generated-by: OpenShiftWebConsole
      creationTimestamp: null
      labels:
        app: hiadmin
        deploymentconfig: hiadmin
    spec:
      containers:
        - env:
            - name: APP_PROFILES_ACTIVE
              value: dev
            - name: SCM_URL
              value: 'http://gitlab.vpclub.cn:8022'
            - name: TZ
              value: Asia/Shanghai
            - name: DOCKER_API_VERSION
              value: '1.24'
            - name: KUBE_WATCH_TIMEOUT
              value: '20'
            - name: API_VERSION
              value: /api/v3
          image: >-
            docker-registry.default.svc:5000/demo/hiadmin@sha256:851f49987fbbf469cbe87c6c26160a1ce7cb1ce5bdb6b1b7d3795127b8a44436
          imagePullPolicy: Always
          name: hiadmin
          ports:
            - containerPort: 7575
              protocol: TCP
            - containerPort: 8080
              protocol: TCP
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /var/lib/docker
              name: volume1
            - mountPath: /var/run/docker.sock
              name: volume2
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
        - hostPath:
            path: /var/lib/docker
            type: ''
          name: volume1
        - hostPath:
            path: /var/run/docker.sock
            type: ''
          name: volume2
  test: false
  triggers:
    - type: ConfigChange
    - imageChangeParams:
        automatic: true
        containerNames:
          - hiadmin
        from:
          kind: ImageStreamTag
          name: 'hiadmin:v1'
          namespace: demo
        lastTriggeredImage: >-
          docker-registry.default.svc:5000/demo/hiadmin@sha256:851f49987fbbf469cbe87c6c26160a1ce7cb1ce5bdb6b1b7d3795127b8a44436
      type: ImageChange
status:
  availableReplicas: 1
  conditions:
    - lastTransitionTime: '2018-12-07T07:12:15Z'
      lastUpdateTime: '2018-12-07T07:12:15Z'
      message: Deployment config has minimum availability.
      status: 'True'
      type: Available
    - lastTransitionTime: '2018-12-07T16:14:31Z'
      lastUpdateTime: '2018-12-07T16:14:32Z'
      message: replication controller "hiadmin-147" successfully rolled out
      reason: NewReplicationControllerAvailable
      status: 'True'
      type: Progressing
  details:
    causes:
      - imageTrigger:
          from:
            kind: DockerImage
            name: >-
              docker-registry.default.svc:5000/demo/hiadmin@sha256:851f49987fbbf469cbe87c6c26160a1ce7cb1ce5bdb6b1b7d3795127b8a44436
        type: ImageChange
    message: image change
  latestVersion: 147
  observedGeneration: 188
  readyReplicas: 1
  replicas: 1
  unavailableReplicas: 0
  updatedReplicas: 1
```

yml中的namespace 请自行更改，在这里需要讲一下env

名称|值|类型|含义
:---:|:---:|:---:|:---:
|APP_PROFILES_ACTIVE|dev|string|所使用的环境|
|SCM_URL|https://gitlab.cn|string|需要授权的gitlab_url|
|API_VERSION|/api/v3|string|所需的gitlab api version |
|TZ|Asia/Shanghai|string|时区|
|KUBE_WATCH_TIMEOUT|20|int  /min|kube api watch 时长默认20 分|
|DOCKER_API_VERSION|1.24|string|docker api 的版本|

### node镜像制作

...

## CRD 资源

console使用的是`kubenetes` 的资源 我们主要创建了 `build` `buildconfig` `deployment` `deploymentconfig` `gatewayconfig` `pipeline` `pipelineconfig` `serviceconfig`
`sourceconfig` `test` `testconfig` 等资源。
如果你需要在集群内创建资源则。

```bash
cd pkg/crd/templates
## 如果使用的是openshift
oc apply -f .

## 如果使用的卡K8S
kubectl apply -f .
```

