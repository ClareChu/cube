# 接口文档

## 环境

dev环境: http://dev.apps.cloud2go.cn/cube/manager

## 更换pod的服本数量

`Request Method` : PUT

`Request Url`    : {{HOST}}/replicas

`Headers`        :  Content-Type:application/json


请求参数

 |字段名|注释|类型|描述|
 |:---:|:---:|:---:|:---:|
 |id|id|int64|| 
 |name|名称|string||
 |namespace|组名|string||
 |version|版本|string|默认版本v1|
 |replicas|副本数|int32||
 |url|回调URL|string||
  
 请求参数试例

 ```json
{
    "id": 1,
	"name":"293889241971163714-353535033",
    "version":"v1",
	"namespace":"testadm",
	"replicas":1,
    "url": "http://www.baidu.com"
}

```

`响应参数`

 |字段名|注释|类型|描述|
 |:---:|:---:|:---:|:---:|
 |code|状态码|string||
 |message|消息|string||
 |data|data|interface||

`成功响应参数试例`

```json
{
    "code": 200,
    "message": "success",
    "data": null
}

```

`错误响应参数试例`

```json
{
    "code": 500,
    "message": "deployments.extensions \"adsfsdf-205004-v11\" not found",
    "data": null
}
```
