package constant

import "errors"

var (
	// NotLoginError 未登录越权访问错误
	NotLoginError = errors.New("not logged in, please log in before accessing")

	// LoginTokenExpiredError Token 已过期
	LoginTokenExpiredError = errors.New("the identity has expired. please login again")

	// LoginIncorrectAccountOrPasswordError 登录失败，账号或密码错误
	LoginIncorrectAccountOrPasswordError = errors.New("login failed, incorrect account or password")
)
