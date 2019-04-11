/**
 * global conf
 * ENV:
 *   WX_TOKEN               --- 微信服务号token
 *   WX_AES_KEY             --- 密文传输AES key
 *   WX_APPID               --- APP ID
 *   WX_APP_SECRET          --- APP Secret
 *   TOKEN_PATH             --- 缓存access token的文件夹
 *   LISTEN_PORT            --- 侦听端口号
 *   LISTEN_HOST            --- 侦听IP地址
 *   WORKER_NUM             --- 并非线程数
 *   SERVICE_PATH           --- 给微信服务提供的访问路径，在微信服务号管理端设置，只需要path部分，如 "/wx"
 *   REDIRECT_PATH          --- 微信网页授权的统一路口，只需要path部分，如 "/redirect"
 *   SERVICE_TIMEOUT        --- 微信请求参数时间戳超时时间，单位秒，如果<=0则不做超时处理
 *   TZ                     --- 时区名称"Asia/Shanghai"
 * Rosbit Xu
 */
package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/rosbit/go-wx-api/conf"
	"time"
)

var (
	ListenHost = ""
	ListenPort = 7080
	TokenStorePath string
	WorkerNum = 5
	WxService string
	WxRedirect string
	WxTimeout = 0
	Loc = time.FixedZone("UTC+8", 8*60*60)
)

func getEnv(name string, result *string, must bool) error {
	s := os.Getenv(name)
	if s == "" {
		if must {
			return fmt.Errorf("env \"%s\" not set", name)
		}
	}
	*result = s
	return nil
}

func CheckGlobalConf() error {
	var err error
	if err = getEnv("WX_TOKEN", &wxconf.WxParams.Token, true); err != nil {
		return err
	}
	var aesKey string
	getEnv("WX_AES_KEY", &aesKey, false)
	if aesKey != "" {
		if err = wxconf.SetAesKey(aesKey); err != nil {
			return err
		}
	}
	if err = getEnv("WX_APP_ID", &wxconf.WxParams.AppId, true); err != nil {
		return err
	}
	if err = getEnv("WX_APP_SECRET", &wxconf.WxParams.AppSecret, true); err != nil {
		return err
	}
	if err = getEnv("SERVICE_PATH", &WxService, true); err != nil {
		return err
	}
	if WxService[0] != '/' {
		WxService = fmt.Sprintf("/%s", WxService)
	}
	if err = getEnv("REDIRECT_PATH", &WxRedirect, true); err != nil {
		return err
	}
	if WxRedirect[0] != '/' {
		WxRedirect = fmt.Sprintf("/%s", WxRedirect)
	}

	if err = getEnv("TOKEN_PATH", &wxconf.TokenStorePath, true); err != nil {
		return err
	}

	getEnv("LISTEN_HOST", &ListenHost, false)
	var p string
	getEnv("LISTEN_PORT", &p, false)
	if p != "" {
		port,_ := strconv.Atoi(p)
		if port > 0 {
			ListenPort = port
		}
	}
	getEnv("WORKER_NUM", &p, false)
	if p != "" {
		workerNum, _ := strconv.Atoi(p)
		if workerNum > 0 {
			WorkerNum = workerNum
		}
	}
	getEnv("SERVICE_TIMEOUT", &p, false)
	if p != "" {
		to, _ := strconv.Atoi(p)
		if to > 0 {
			WxTimeout = to
		}
	}
	getEnv("TZ", &p, false)
	if p != "" {
		if loc, err := time.LoadLocation(p); err == nil {
			Loc = loc
		}
	}
	return nil
}

func DumpConf() {
	fmt.Printf("listen host: %s\n", ListenHost)
	fmt.Printf("listen port: %d\n", ListenPort)
	fmt.Printf("token store path: %s\n", wxconf.TokenStorePath)
	fmt.Printf("task handler count: %d\n", WorkerNum)
	fmt.Printf("Wx Params: %v\n", wxconf.WxParams)
	fmt.Printf("Wx Service: %s\n", WxService)
	fmt.Printf("Redirect Service: %s\n", WxRedirect)
	fmt.Printf("Service Timeout: %d\n", WxTimeout)
	fmt.Printf("TZ time location: %v\n", Loc)
}
