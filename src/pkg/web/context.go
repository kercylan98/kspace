package web

import "github.com/gin-gonic/gin"

// Context 对于 gin.Context 的变种
type Context struct {
	*gin.Context
}
