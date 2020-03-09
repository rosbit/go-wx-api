/**
 * 缺省的微信网页授权处理
 */
package wxauth

type RedirectHandlerWithoutAppId func(openId, state string) (string, map[string]string, string, error)

// 根据服务号菜单state做跳转的实现，缺省实现可以被RegisterRedictHandler覆盖
var HandleRedirect = func(openId, state string) (string, map[string]string, string, error) {
	return "success", nil, "", nil
}

// ---------------- 支持多服务号的实现 ------------------
/**
 * [函数签名]根据服务号菜单state做跳转
 * @param appId   服务号的appId，用于区分服务号
 * @param openId  订阅用户的openId
 * @param state   微信网页授权中的参数，用来标识某个菜单
 * @return
 *   c    需要显示服务号对话框中的内容
 *   h    需要在微信内嵌浏览器中设置的header信息，包括Cookie
 *   r    需要通过302跳转的URL。如果r不是空串，c的内容被忽略
 *   err  如果没有错误返回nil，非nil表示错误
 */
type RedirectHandler func(appId, openId, state string) (c string, h map[string]string, r string, err error)

func ToAppIdRedirectHandler(handler RedirectHandlerWithoutAppId) RedirectHandler {
	return func(appId,openId,state string)(c string,h map[string]string,r string,err error) {
		c,h,r,err = handler(openId, state)
		return
	}
}

// 注册转发HTTP(s) URL，该URL将全权决定网页授权的处理。如果该URL存在，优先级要"高于"RegisterRedictHandler()注册函数。
// 参数JSON: {"appId": "xxx", "openId": "xxx", "state": "state"}
// 该URL的以POST形式接收参数，而且会得到所有的HTTP头信息，可以设置任何的响应头信息，响应结果直接显示在公众号浏览器中
// 响应时间要控制好，避免微信服务超时
func (p *WxAppIdAuthHandler) RegisterRedirectUrl(redirectUrl string) {
	p.redirectUrl = redirectUrl
}

// 注册网页授权处理函数，该函数优先级"低于"RegisterRedirectUrl()注册的URL
func (p *WxAppIdAuthHandler) RegisterRedictHandler(handler RedirectHandler) {
	if handler != nil {
		p.redirectHandler = handler
	}
}
