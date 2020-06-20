# `rpc option` 测试用例

## option alias 

- 在 `rpc` 的 `body` 中定义 `option`，形如

```go
service helloworld_svr {
    rpc Hello(HelloReq) returns(HelloRsp) {
        option(gorpc.alias) = "/api/v1/helloworld";
    };
}
```

- `gorpc.alias` 重命名 `rpc` 的方法名。

## option swagger 

- 在 `rpc` 的 `body` 中定义 option，形如

```go
service helloworld_svr {
    rpc Hello(HelloReq) returns(HelloRsp) {
        option(gorpc.swagger) = {
            title : "你好世界"
            method: "get"
            description:
                "入参：msg\n"
                "作用：用于演示 helloword\n"
        };
    };
}
```

- `gorpc.swagger` 的 `title` 为该 `rpc` 的方法名，`method` 为 `http` 的请求方法（如果该接口用于 `http`，
由于 `swagger-ui` 会识别一个 `method`，如果该字段不填，默认为 `post`），`description` 用于描述此接口。

## 使用方法

1 pb 中加入 `option` 定义。

2 命令参数中加入 `swagger` 或 `alias`，如：

```shell script
gorpc create -p=./cmd/testcase.option/helloworld.proto -o=project --swagger --force -v --alias
```

3 在当前目录下会生成 `apidocs.swagger.json` 

4 下载 `swagger-ui` (https://github.com/swagger-api/swagger-ui)

5 进入到仓库下的 `dist` 目录，将 `apidocs.swagger.json` 拷贝至此，并修改 `index.html` 文件中的 `url` 为 `apidocs.swagger.json`。

6 `npm install -g http-server`，直接运行 `http-server` 后可以通过 bash 显示的 url 对 swagger 页面进行访问。

