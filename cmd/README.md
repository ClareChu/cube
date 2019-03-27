# cube 命令行客户端 cube/client

cube/client 采用Go语言开发，能跨平台运行在 Windows, Mac, Linux.

在你编译前，请确保你安装了Go 1.11 或以上版本。

## 编译及安装方法

```bash

git clone https://github.com/hidevopsio/cube.git

cd cube/client

go build -o cube

cp cube $GOPATH/bin/

```

## 使用方法

cube -h

```bash
  info        Display client configuration information
  login       Log in to a cube server
  logs        Fetch the logs of a container
  run         Run a command in a new pipeline
  set         Set client configuration information
```

```bash
cube set --name={name} --token={token} --host={host_url}
cube info
cube login
cube run --project={namspaces} --app={app_name} --sourcecode={java|nodejs|go}
cube logs --project={namspaces} --app={app_name} --sourcecode={java|nodejs|go}
```

examples:
```
cube set --host=http://cube-server.examples.com
cube login
cube run --project=demo --app=hello-world --sourcecode=java
cube logs --project=demo --app=hello-world --sourcecode=java
```
