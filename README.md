# 微信公众号API封装

go-wx-api是对微信公众号API的封装，可以当作SDK使用

## 下载方法

 1. 命令行进入任一文件夹
 2. 执行命令
    ```bash
    GOPATH=`pwd` go get -u github.com/rosbit/go-wx-api
    ```

## 编译例子
 1. 在go-wx-api所在文件夹下
 2. 执行命令
    ```bash
    GOPATH=`pwd` go get -u github.com/rosbit/go-wx-api/samples/wx-echo-server
    ```
    命令执行成功后，会有一个可执行程序`bin/wx-echo-server`，该程序可以直接用于配置微信公众号**服务器配置**，并可以对公众号对话框输入做回声应答
 4. 执行命令
    ```bash
    GOPATH=`pwd` go install github.com/rosbit/go-wx-api/samples/wx-server
    ```
    可以得到`bin/wx-server`，一个可以处理菜单的公众号服务

## 使用方法

以下是一个简单的例子，详细的例子参考`samples/wx-echo-server`和`samples/wx-server`

```go
package main

import (
	"github.com/rosbit/go-wx-api/conf"
	"github.com/rosbit/go-wx-api"
	"net/http"
	"fmt"
)

const (
	token     = "微信公众号的token"
	appId     = "微信公众号appId"
	appSecret = "微信公众号的secret"
	aesKey    = "" //安全模式 使用的AESKey，如果是 明文传输，该串为空
	
	listenPort = 7070   // 服务侦听的端口号，请根据微信公众号管理端的服务器配置正确设置
	service    = "/wx"  // 微信公众号管理端服务器配置中URL的路径部分
)

func main() {
	wxconf.WxParams = wxconf.WxParamsT{Token:token, AppId:appId, AppSecret:appSecret}
	if aesKey != "" {
		if err := wxconf.SetAesKey(aesKey); err != nil {
			fmt.Printf("invalid aesKey: %v\n", err)
			return
		}
	}

	http.HandleFunc(service, wxapi.Echo)     // 用于配置
	http.HandleFunc(service, wxapi.Request)  // 用于实际执行公众号请求，和wxapi.Echo只能使用一个。
	                                         // 可以使用高级路由功能同时设置，参考 github.com/rosbit/go-wx-api/samples/wx-echo-server
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
```

## 微信工具类程序
 1. 在go-wx-api所在文件夹下
 2. 执行命令
     ```bash
     GOPATH=`pwd` go install github.com/rosbit/go-wx-api/tools/create-wx-menu
     ```
     可以根据配置文件创建公众号的自定义菜单，5分钟内生效
