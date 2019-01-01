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

如果使用GO版本在 `go1.11` 以下 请自行安装dep包管理工具。
在项目下执行

```bash
dep ensure -v
```

如果使用1.11 以上版本最则使用go modules
该项目下分为`console` 和 `node` 俩个项目, console 主要负责CI/CD 调度 启动 任务分发 等作用， node 主要负责 代码编译 代码测试 镜像的制作等主要工作。

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

docker build -t docker-registry-default.app.example.io/demo/hiadmin:v1 .

docker push docker-registry-default.app.example.io/demo/hiadmin:v1

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
              value: 'http://gitlab.example.cn:8022'
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
进入node目录执行make命令即可

make

如需变更版本和镜像名称修改node目录下Makefile中变量

imageName = {registries}/{group}/{image_name}:tag
```
tag := 1.1.8
registries := docker.example.cn
group := hidevopsio
image_name := hinode-java-jar
```

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

### 1. buildConfig

 `buildConfig` 资源主要起到的作用是编译 打包代码 代码检测 镜像的编译等作用，以下是讲解buildconfig字段的含义。

```yaml
apiVersion: mio.io/v1alpha1
kind: BuildConfig
metadata:
  name: java
  namespace: templates
spec:
  app: ""
  baseImage: docker.example.cn/hidevopsio/hinode-java-jar:1.1.12
  cloneConfig:
    branch: master
    depth: 1
    dstDir: /opt/app-root/src/example
    password: ""
    url: https://gitlab.example.cn
    username: ""
  cloneType: ""
  codeType: java
  compileCmd:
  - commandName: pwd
  - Script: |-
      mvn clean package -U -Dmaven.test.skip=true -Djava.net.preferIPv4Stack=true
      if [[ $? == 0 ]]; then
        echo 'Build Successful.'
      else
        echo 'Build Failed!'
        exit 1
      fi
    execType: script
  - commandName: ls
  - params： -a
  deployData:
    envs:
      CODE_TYPE: java
      DOCKER_API_VERSION: "1.24"
      MAVEN_HOST: http://nexus.example.cn
      MAVEN_MIRROR_URL: http://nexus.example.cn/repository/maven-public/
      NODE_NAME: node05.example.io
    hostPathVolume:
      /var/lib/docker: /var/lib/docker
      /var/run/docker.sock: /var/run/docker.sock
    ports:
    - 8080
    - 7575
    replicas: 1
  dockerAuthConfig:
    password: Harbor12345
    username: admin
  dockerFile:
  - FROM docker.example.cn/hidevopsio/java:8-jre-alpine
  - ENV  TZ="Asia/Shanghai"
  - ENV  APP_OPTIONS="-Xms128m -Xmx512m -Xss512k"
  - ENV   APP_OPTIONS="-Xms128m -Xmx512m -Xss512k"
  - COPY app.jar /root
  - EXPOSE 8080
  - EXPOSE 7575
  - ENTRYPOINT ["sh","-c","java -jar /root/app.jar $APP_OPTIONS"]
  dockerRegistry: harbor.example.io
  nodeService: ""
  tasks:
  - name: createService
  - name: deployNode
  - name: clone
  - name: compile
  - name: buildImage
  - name: pushImage
status:
  lastVersion: 1
```

字段|类型|含义
:---:|:---:|:---:
|spec.baseImage|string|mio下node项目制作的镜像名称|
|spec.cloneType|string|代码克隆模式|
|spec.codeType|string|代码语言类型|
|spec.dockerRegistry|string|镜像仓库地址|
|spec.tasks|string|pipeline任务序列|
|---|---|---|
|spec.cloneConfig|object|存放克隆代码所需的一些信息|
|spec.cloneConfig.url|string|设置克隆代码的http地址|
|spec.cloneConfig.bransh|string|配置即将克隆代码的分支|
|spec.cloneConfig.depth|string|配置代码克隆的深度|
|spec.cloneConfig.dstDir|string|代码克隆在容器内部的目标地址|
|spec.cloneConfig.username|string|配置克隆代码的鉴权用户名信息|
|spec.cloneConfig.password|string|配置克隆代码的鉴权用密码信息|
|---|---|---|
|spec.compileCmd|list|用来存放一个shell命令组，主要用来代码编译|
|spec.compileCmd`[`0`]`.execType|string|用来选择是否执行script字段中的脚本内容|
|spec.compileCmd`[`0`]`.script|string|execType选择了script后会使用此字段中的内容生成一个脚本并执行|
|spec.compileCmd`[`0`]`.commandName|string|command名称，使用时不能设置execType为script|
|spec.compileCmd`[`0`]`.params|string|command参数和commandName配合使用|
|---|---|---|
|spec.deployData|object|用于node发布时的一些配置数据|
|spec.deployData.envs|list|以key，value形式存放环境变量的一个数组|
|spec.deployData.hostPathVolume|list|以key，value形式存放hostPathVolume挂载的一个数组|
|spec.deployData.ports|list|用来存放node服务端口暴露的配置字段|
|---|---|---|
|spec.dockerAuthConfig|object|用于配置docker权限的配置信息|
|spec.dockerAuthConfig.username|string|用户名|
|spec.dockerAuthConfig.password|string|密码|
|---|---|---|
|spec.dockerFile|list|一个以行为单位存放Dockerfile语句的数组|

### 2. deploymentConfig

`deploymentConfig` 主要是针对镜像deploy

```yaml
apiVersion: mio.io/v1alpha1
kind: DeploymentConfig
metadata:
  name: java
  namespace: templates
spec:
  dockerRegistry: docker-registry.default.svc:5000
  env:
  - name: starter
    value: jav -jar
  - name: TZ
    value: Asia/Shanghai
  - name: APP_OPTIONS
    value: -Xms128m -Xmx512m -Xss512k
  - name: SPRING_PROFILES_ACTIVE
    value: dev
  envType:
  - remoteDeploy
  - deploy
  fromRegistry: docker-registry-default.app.example.io
  port:
  - containerPort: 8080
    name: tcp-8080
    protocol: TCP
  profile: dev
status:
  lastVersion: 1
```

字段|类型|含义|
:---:|:---:|:---:
|dockerRegistry|string|推镜像的仓库地址|
|env|[]string|pod 环境变量|
|envType|string| 执行步骤|
|fromRegistry|string|拉取镜像仓库地址|
|port|---|对外端口|
|profile|string|环境变量|

### 3. serviceConfig

```yaml
apiVersion: mio.io/v1alpha1
kind: ServiceConfig
metadata:
  generation: 0
  name: java
  namespace: templates
spec:
  ports:
  - name: http-8080
    port: 8080
    protocol: TCP
    targetPort: 8080
  - name: http-7575
    port: 7575
    protocol: TCP
    targetPort: 7575
status:
  metadata: {}
```
