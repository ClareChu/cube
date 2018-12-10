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

如果使用GO版本在 `go1.11` 下 请自行安装dep包管理工具。
在项目下执行

```bash
dep ensure -v
```

如果使用1.11 以上版本最则使用go modules
该项目下分为`console` 和 `node` 俩个项目, console 主要负责CICD 调度 启动 任务分发 等作用， node 主要负责 代码编译 代码测试 镜像的制作等主要工作。

### console镜像制作 与部署

创建namespace

```bash
oc create namespace hidevopsio

### 给namespace 授权
oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:hidevopsio:default

oc adm policy add-scc-to-user privileged -z default -n hidevopsio

### 默认情况openshift 会预先分配UID 不能以任何用户身份运行，并阻止容器运行, 修改restricted

oc edit scc restricted

# 更改runAsUser.Type为RunAsAny。
# 确保allowPrivilegedContainer设置为false。

# 要修改群集，以便它不预先分配UID并且不允许容器以root身份运行：
# 更改runAsUser.Type为MustRunAsNonRoot。

# 使用hostPath卷插件
#  要放松群集中的安全性，以便允许群集使用 hostPath卷插件而不授予每个人访问特权 SCC 的权限：
# 编辑受限制的 SCC：
oc edit scc restricted
#  添加 allowHostDirVolumePlugin: true。
```

```bash
cd console
### 执行build.sh
./build.sh
###脚本大致内容

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
  labels:
    app: hiadmin
  name: hiadmin
  namespace: hidevopsio
  resourceVersion: '139885254'
  selfLink: /apis/apps.openshift.io/v1/namespaces/demo/deploymentconfigs/hiadmin
  uid: d2bcdad9-d4ff-11e8-bb8c-005056935c80
spec:
  replicas: 1
  selector:
    app: hiadmin
    deploymentconfig: hiadmin
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

