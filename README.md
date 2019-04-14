# 微信公众号开发SDK

go-wx-api是对微信公众号API的封装，可以当作SDK使用

## 编译例子
 1. 该函数包已经使用go modules发布，需要golang 1.11.x及以上版本
 1. 请参考[go-wx-apps](https://github.com/rosbit/go-wx-apps)，那里包含了例程和工具程序

## 使用方法

以下是一个简单的例子，用于说明使用go-wx-api的主要执行步骤。更详细的例子参考[go-wx-apps](https://github.com/rosbit/go-wx-apps)

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

	workerNum = 3 // 处理请求的并发数
)

func main() {
	// 步骤1. 设置配置参数
	wxconf.WxParams = wxconf.WxParamsT{Token:token, AppId:appId, AppSecret:appSecret}
	if aesKey != "" {
		if err := wxconf.SetAesKey(aesKey); err != nil {
			fmt.Printf("invalid aesKey: %v\n", err)
			return
		}
	}

	// 步骤2. 初始化SDK
	wxapi.InitWxAPI(workerNum, os.Stdout)

	// 步骤2.5 设置签名验证的中间件。由于net/http不支持中间件，省去该步骤
	// signatureChecker := wxapi.NewWxSignatureChecker(wxconf.WxParams.Token, 0, []string{service})
	// <middleWareContainer>.Use(signatureChecker)

	// 步骤3. 设置http路由，启动http服务
	http.HandleFunc(service, wxapi.Echo)     // 用于配置
	http.HandleFunc(service, wxapi.Request)  // 用于实际执行公众号请求，和wxapi.Echo只能使用一个。
	                                         // 可以使用高级路由功能同时设置，参考 github.com/rosbit/go-wx-api/samples/wx-echo-server
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
```

## 其它
 1. 该函数包可以处理文本消息、用户关注/取消关注事件、菜单点击事件
 2. 其它消息、事件可以根据需要扩充
