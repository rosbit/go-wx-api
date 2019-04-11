/**
 * 缺省的微信网页授权处理
 */
package wxauth

// 根据服务号菜单state做跳转的实现，缺省实现可以被RegisterRedictHandler覆盖
var HandleRedirect = func(openId, state string) (string, map[string]string, string, error) {
	return "success", nil, "", nil
}

