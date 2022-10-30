# 微信公众号开发SDK

go-wx-api是对微信公众号API的封装，可以当作SDK使用，主要特点:

 - 对常用消息、事件的接收和回复做了封装，已经无需了解相关的公众号开发文档；见使用方法1；
 - 提供了缺省的消息、事件处理方法，可以根据实际需求覆盖相关实现；见使用方法2；
 - 使用`go-wx-api`可以同时支持**多个公众号**；见使用方法3;
 - 菜单点击后允许进入自己的页面，go-wx-api提供了统一的入口进入菜单处理，获取用户openId，
   根据state参数进入具体业务处理。可以通过提供菜单跳转处理器实现具体业务；见使用方法1;
 - 已经支持视频号小店的推送接口、查看订单/售后单详情的接口
 - [go-wx-gateway](https://github.com/rosbit/go-wx-gateway)是使用go-wx-api实现的微信公众号网关服务，
   通过`go-wx-gateway`，就可以把微信公众号服务的开发转化为普通web服务开发。

## 编译例子
 1. 该函数包已经使用go modules发布，需要golang 1.11.x及以上版本
 1. 请参考[go-wx-apps](https://github.com/rosbit/go-wx-apps)，那里包含了例程和工具程序

## 使用方法1: (实现消息处理器、菜单跳转处理器的例子)

go-wx-api已经对公众号常用的消息(文本框架输入、发语音等)、事件(用户关注、点击菜单等)做了提取和封装处理。缺省处理对消息、事件做了简单的应答处理，缺省处理除了能给公众号管理后台做开发者设置功能外，没有实际意义。具体业务可以根据需要对go-wx-api的消息处理器接口进行实现：

 1. 消息处理器的定义和实现

    ```go
    import (
        "github.com/rosbit/go-wx-api/v2/msg"
        "fmt"
    )

    // 消息处理器定义
    type YourMsgHandler struct {
        wxmsg.WxMsgHandlerAdapter  // 包含了所有消息、事件的缺省实现
    }

    // 接口定义见 wxmsg.WxMsgHandler，根据需要选择实现其中的方法

    // 文本消息处理
    func (h *YourMsgHandler) HandleTextMsg(textMsg *wxmsg.TextMsg) wxmsg.ReplyMsg {
        return NewReplyTextMsg(textMsg.FromUserName, textMsg.ToUserName, fmt.Sprintf("收到了消息:%s", textMsg.Content))
    }

    // 用户关注公众号处理
    func (h *YourMsgHandler) HandleSubscribeEvent(subscribeEvent *wxmsg.SubscribeEvent) wxmsg.ReplyMsg {
        return wxmsg.NewReplyTextMsg(subscribeEvent.FromUserName, subscribeEvent.ToUserName, "welcome")
    }
    ```

 1. 注册消息处理器
    - 单一公众号注册方法见方法2
    - 多公众号注册方法见方法3

## 使用方法2: (单一公众号服务)

以下是一个简单的例子，用于说明使用go-wx-api的主要执行步骤。更详细的例子参考[go-wx-apps](https://github.com/rosbit/go-wx-apps)

```go
package main

import (
	"github.com/rosbit/go-wx-api/v2/conf"
	"github.com/rosbit/go-wx-api/v2/msg"
	"github.com/rosbit/go-wx-api/v2"
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
	name = "test" // 服务的名字
)

func main() {
	// 步骤1. 初始化accessToken缓存路径
	wxapi.InitWx("/tmp")

	// 步骤2. 设置配置参数
	if err := wxapi.SetWxParams(name, token, appId, appSecret, aesKey); err != nil {
		fmt.Printf("failed to set params: %v\n", err)
		return
	}

	// 步骤3: 设置签名验证的中间件。由于net/http不支持中间件，省去该步骤
	// signatureChecker := wxapi.NewWxSignatureChecker(token, 0, []string{service})
	// <middleWareContainer>.Use(signatureChecker)

	// 步骤4. 设置http路由，启动http服务
	// 用于公众号配置
	http.HandleFunc(service, wxapi.CreateEcho(token))

	// 用于实际执行公众号请求，和wxapi.CreateEcho只能使用一个。
	// 可以使用支持高级路由功能的web框架同时设置
	http.HandleFunc(service, wxapi.CreateMsgHandler(name, workerNum, wxmsg.MsgHandler))  // 使用缺省的消息处理器
	http.HandleFunc(service, wxapi.CreateMsgHandler(name, workerNum, &YourMsgHandler{})) // 使用自定义的消息处理器
	                                         
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
```

## 使用方法3: (多个公众号服务)

以下代码仅仅为同时启用公众号的示例:

```go
package main

import (
	"github.com/rosbit/go-wx-api/v2/conf"
	"github.com/rosbit/go-wx-api/v2"
	"net/http"
	"fmt"
)

type WxConf struct {
	name string
	token string
	appId string
	appSecret string
	aesKey string
	workerNum int
	service string
	rd string
	redirectUrl string
}

var (
	listenPort = 7070   // 服务侦听的端口号，请根据微信公众号管理端的服务器配置正确设置
	wxServices = []WxConf{
		WxConf{
			name: "wx1",
			token: "微信公众号1的token",
			appId: "微信公众号1的appId",
			appSecret: "微信公众号的1secret",
			aesKey: "",      // 安全模式 使用的AESKey，如果是 明文传输，该串为空
			workerNum: 3,    // 处理请求的并发数
			service: "/wx1", // 微信公众号管理端服务器配置中URL的路径部分
			rd: "/rd1", // 网页授权接受地址
			redirectUrl: "http://localhost:port/redirect", // 接受网页授权转发的其它服务
		},
		WxConf{
			name: "wx2",
			token: "微信公众号2的token",
			appId: "微信公众号2的appId",
			appSecret: "微信公众号2的secret",
			aesKey: "",      // 安全模式 使用的AESKey，如果是 明文传输，该串为空
			workerNum: 3,    // 处理请求的并发数
			service: "/wx2", // 微信公众号管理端服务器配置中URL的路径部分
			rd: "/rd2", // 网页授权接受地址
			redirectUrl: "http://localhost:port/redirect", // 接受网页授权转发的其它服务
		},
		// 其它服务号
	}
)

func main() {
	// 步骤1. 初始化accessToken缓存路径
	wxapi.InitWx("/tmp")

	// 对于每一个公众号执行
	for _, conf := range wxServices {
		// 步骤2. 设置配置参数
		if err := wxapi.SetWxParams(conf.name, conf.token, conf.appId, conf.appSecret, conf.aesKey); err != nil {
			fmt.Printf("failed to set params: %v\n", err)
			return
		}

		// 步骤2.5 设置签名验证的中间件。由于net/http不支持中间件，省去该步骤
		// signatureChecker := wxapi.NewWxSignatureChecker(conf.token, 0, []string{conf.service})
		// <middleWareContainer>.Use(signatureChecker)

		// 步骤3. 设置http路由，启动http服务
		http.HandleFunc(conf.service, wxapi.CreateEcho(conf.token))     // 用于配置

		// 用于实际执行公众号请求，和wxapi.CreateEcho只能使用一个。
		// 可以使用支持高级路由功能的web框架同时设置。使用POST路由
		http.HandleFunc(conf.service, wxapi.CreateMsgHandler(conf.name, conf.workerNum, wxmsg.MsgHandler))  // 使用缺省的消息处理器
		http.HandleFunc(conf.service, wxapi.CreateMsgHandler(conf.name, conf.workerNum, &YourMsgHandler{})) // 使用自定义的消息处理器

		// 菜单跳转全权交给另外一个URL处理，用GET路由
		// redirectURL接收POST请求，POST body是一个JSON: {"appId":"xxx", "openId", "xxx", "state": "xxx", "userInfo": {}}
		// 它可以随意处理HTTP请求、输出HTTP响应，响应结果直接返回公众号浏览器
		http.HandleFunc(conf.rd, wxapi.CreateOAuth2Redirector(conf.name, conf.workerNum, conf.redirectUrl))
	}

	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}

```

## 其它
 1. 该函数包可以处理文本消息、用户关注/取消关注事件、菜单点击事件
 2. 其它消息、事件可以根据需要扩充
