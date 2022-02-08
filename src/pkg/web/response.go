package web

import (
	"encoding/json"
	"errors"
)

// Response 响应体结构
//
// 交由 Server 在接口处理完毕后对结果进行统一响应，如果提前通过 gin.ResponseWriter 进行响应后，将不会再重复响应
type Response struct {
	Successful bool           `json:"successful"` // 请求是否成功
	Data       any            `json:"data"`       // 请求包含的数据
	Error      *ResponseError `json:"error"`      // 请求包含的错误信息

	noWriter bool // 指定无需响应
}

// ResponseError 响应体错误信息结构
type ResponseError struct {
	// 指针：对于 Response 在还未完全构建完成 ResponseError 的时候，如果存入了 Response 的值，那么此时 Response.Err 还未生
	// 效，所以需要存入指针确保最后执行 Throw 时能够正确返回结果
	response    *Response
	Code        int    `json:"code"`        // 错误代码
	Route       string `json:"route"`       // 发生错误的路由（该由程序赋予）
	Explain     string `json:"explain"`     // 来自系统的错误说明
	Suggestions string `json:"suggestions"` // 解决方案提议
	Amicable    string `json:"amicable"`    // 友好的可以供给客户查阅的提示
}

// NoReply 设置后将不会对请求结果进行响应
func (slf *Response) NoReply() Response {
	slf.noWriter = true
	return *slf
}

// Pass 表示已经成功执行完毕，并且需要携带一些数据（data）进行返回
//
// data 传入单个和多个之间存在一些区别
//
// 单个：
//
// {data: {...}}
//
// 多个：
//
// {data: [{...}, {...}]}
func (slf *Response) Pass(data ...any) Response {
	var dataCount = len(data)
	if dataCount == 1 {
		slf.Data = data[0]
	} else if dataCount > 1 {
		slf.Data = data
	}
	slf.Successful = true
	return *slf
}

// ErrJSON 定义JSON格式内容的错误（err）已经发生，并带有默认的错误代码：500（可以被 http.Response 的 StatusCode 改变）
//
// 如果非 JSON 格式错误，将转到默认的 Response.Err 处理
func (slf *Response) ErrJSON(err error) *ResponseError {
	var data = map[string]interface{}{}
	if unmarshalErr := json.Unmarshal([]byte(err.Error()), &data); unmarshalErr != nil {
		return slf.Err(err)
	}
	var errName, errDescribe interface{}
	var exist bool
	if errName, exist = data["error"]; !exist {
		return slf.Err(err)
	}
	if errDescribe, exist = data["error_description"]; !exist {
		return slf.Err(err)
	}
	var responseErr *ResponseError
	switch v := errName.(type) {
	case string:
		responseErr = slf.Err(errors.New(v))
	default:
		return slf.Err(err)
	}
	switch v := errDescribe.(type) {
	case string:
		responseErr.MaybeSo(v)
	default:
		return slf.Err(err)
	}
	slf.Error = responseErr
	return slf.Error
}

// Err 定义错误（err）已经发生，并带有默认的错误代码：500（可以被 http.Response 的 StatusCode 改变）
func (slf *Response) Err(err error) *ResponseError {
	var errStr = err.Error()
	slf.Successful = false
	slf.Data = nil
	slf.Error = &ResponseError{
		response:    slf,
		Code:        500,
		Explain:     errStr,
		Suggestions: "no suggestions",
		Amicable:    errStr,
	}
	return slf.Error
}

// ErrWithCode 定义错误（err）已经发生并且声明一个错误代码（code）
func (slf *Response) ErrWithCode(code int, err error) *ResponseError {
	slf.Err(err).Code = code
	return slf.Error
}

// MaybeSo 解决方案建议（suggestions）
func (slf *ResponseError) MaybeSo(suggestions string) *ResponseError {
	slf.Suggestions = suggestions
	return slf
}

// Show 应该反馈给用户查看的内容（amicable）
func (slf *ResponseError) Show(amicable string) *ResponseError {
	slf.Amicable = amicable
	return slf
}

// Throw 抛出错误，应该在错误完全定义完毕的时候使用
func (slf ResponseError) Throw() Response {
	return *slf.response
}
