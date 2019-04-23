# 微信公众号开发SDK

go-wx-api是对微信公众号API的封装，可以当作SDK使用，主要特点:

 - 对常用消息、事件的接收和回复做了封装，已经无需了解相关的公众号开发文档；见使用方法1；
 - 提供了缺省的消息、事件处理方法，可以根据实际需求覆盖相关实现；见使用方法2；
 - 使用go-wx-api可以同时支持多个公众号；见使用方法3;
 - 菜单点击后允许进入自己的页面，go-wx-api提供了统一的入口进入菜单处理，获取用户openId，
   根据state参数进入具体业务处理。可以通过提供菜单跳转处理器实现具体业务；见使用方法1;

## 编译例子
 1. 该函数包已经使用go modules发布，需要golang 1.11.x及以上版本
 1. 请参考[go-wx-apps](https://github.com/rosbit/go-wx-apps)，那里包含了例程和工具程序

## 使用方法1: (实现消息处理器、菜单跳转处理器的例子)

go-wx-api已经对公众号常用的消息(文本框架输入、发语音等)、事件(用户关注、点击菜单等)做了提取和封装处理。缺省处理对消息、事件做了简单的应答处理，缺省处理除了能给公众号管理后台做开发者设置功能外，没有实际意义。具体业务可以根据需要对go-wx-api的消息处理器接口进行实现：

 1. 消息处理器的定义和实现

    ```go
    import (
        "github.com/rosbit/go-wx-api/msg"
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

 1. 菜单跳转处理器实现
    - 单一公众号

       ```go
       /**
        * 据服务号菜单state做跳转
        * @param openId  订阅用户的openId
        * @param state   微信网页授权中的参数，用来标识某个菜单
        * @return
        *   c    需要显示服务号对话框中的内容
        *   h    需要在微信内嵌浏览器中设置的header信息，包括Cookie
        *   r    需要通过302跳转的URL。如果r不是空串，c的内容被忽略
        *   err  如果没有错误返回nil，非nil表示错误
        */
       func handleMenuRedirect(openId, state string) (c string, h map[string]string, r string, err error) {
            r = "http://www.yourhost.com/path/to/service"
            return
       }
       ```

    - 多公众号

       ```go
       /**
        * 据服务号菜单state做跳转
        * @param appId   公众号的appId，用于区分不同的公众号
        * @param openId  订阅用户的openId
        * @param state   微信网页授权中的参数，用来标识某个菜单
        * @return
        *   c    需要显示服务号对话框中的内容
        *   h    需要在微信内嵌浏览器中设置的header信息，包括Cookie
        *   r    需要通过302跳转的URL。如果r不是空串，c的内容被忽略
        *   err  如果没有错误返回nil，非nil表示错误
        */
       func handleMenuRedirect(appId, openId, state string) (c string, h map[string]string, r string, err error) {
            r = "http://www.yourhost.com/path/to/service"
            return
       }
       ```

## 使用方法2: (单一公众号服务)

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
	if err := wxconf.SetParams(token, appId, appSecret, aesKey); err != nil {
		fmt.Printf("failed to set params: %v\n", err)
		return
	}

	// 步骤2. 初始化SDK
	wxapi.InitWxAPI(workerNum, os.Stdout)

	// 注册消息处理器、菜单跳转处理器。如果没有相应的实现，可以注释掉下面两行代码
	wxapi.RegisterWxMsghandler(&YourMsgHandler{})
	wxapi.RegisterRedictHandler(handleMenuRedirect)

	// 步骤2.5 设置签名验证的中间件。由于net/http不支持中间件，省去该步骤
	// signatureChecker := wxapi.NewWxSignatureChecker(wxconf.WxParams.Token, 0, []string{service})
	// <middleWareContainer>.Use(signatureChecker)

	// 步骤3. 设置http路由，启动http服务
	http.HandleFunc(service, wxapi.Echo)     // 用于配置
	http.HandleFunc(service, wxapi.Request)  // 用于实际执行公众号请求，和wxapi.Echo只能使用一个。
	                                         // 可以使用支持高级路由功能的web框架同时设置，参考 github.com/rosbit/go-wx-api/samples/wx-echo-server
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
```

## 使用方法3: (多个公众号服务)

以下代码仅仅为同时启用公众号的示例:

```go
package main

import (
	"github.com/rosbit/go-wx-api/conf"
	"github.com/rosbit/go-wx-api"
	"net/http"
	"fmt"
)

type WxConf struct {
	token string
	appId string
	appSecret string
	aesKey string
	workerNum int
	service string
}

var (
	listenPort = 7070   // 服务侦听的端口号，请根据微信公众号管理端的服务器配置正确设置
	wxServices = []WxConf{
		WxConf{
			token: "微信公众号1的token",
			appId: "微信公众号1的appId",
			appSecret: "微信公众号的1secret",
			aesKey: "",      // 安全模式 使用的AESKey，如果是 明文传输，该串为空
			workerNum: 3,    // 处理请求的并发数
			service: "/wx1", // 微信公众号管理端服务器配置中URL的路径部分
		},
		WxConf{
			token: "微信公众号2的token",
			appId: "微信公众号2的appId",
			appSecret: "微信公众号2的secret",
			aesKey: "",      // 安全模式 使用的AESKey，如果是 明文传输，该串为空
			workerNum: 3,    // 处理请求的并发数
			service: "/wx2", // 微信公众号管理端服务器配置中URL的路径部分
		},
		// 其它服务号
	}
)

func main() {
	// 对于每一个公众号执行
	for _, conf := range wxServices {
		// 步骤1. 设置配置参数
		wxParams, err := wxconf.NewWxParams(conf.token, conf.appId, conf.appSecret, conf.aesKey)
		if err != nil {
			fmt.Printf("failed to set params: %v\n", err)
			return
		}

		// 步骤2. 初始化SDK
		wxService := wxapi.InitWxAPIWithParams(wxParams, conf.workerNum, os.Stdout)

		// 注册消息处理器、菜单跳转处理器。如果没有相应的实现，可以注释掉下面两行代码
		wxService.RegisterWxMsghandler(&YourMsgHandler{})   // 不同的wxService可以有不同的MsgHandler
		wxService.RegisterRedictHandler(handleMenuRedirect) // 不同的wxSercice可以有不同的RedirectHandler

		// 步骤2.5 设置签名验证的中间件。由于net/http不支持中间件，省去该步骤
		// signatureChecker := wxapi.NewWxSignatureChecker(wxParams.Token, 0, []string{conf.service})
		// <middleWareContainer>.Use(signatureChecker)

		// 步骤3. 设置http路由，启动http服务
		http.HandleFunc(conf.service, wxService.Echo)     // 用于配置
		http.HandleFunc(conf.service, wxService.Request)  // 用于实际执行公众号请求，和wxService.Echo只能使用一个。
		                                                  // 可以使用支持高级路由功能的web框架同时设置
	}

	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}

```

## 其它
 1. 该函数包可以处理文本消息、用户关注/取消关注事件、菜单点击事件
 2. 其它消息、事件可以根据需要扩充
