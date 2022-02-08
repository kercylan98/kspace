# 关于 API 路由应该这样定义

```
var accounts = server.Group("/accounts")

accounts.GET(":id", c.ShowAccount)
accounts.GET("", c.ListAccounts)
accounts.POST("", c.AddAccount)
accounts.DELETE(":id", c.DeleteAccount)
accounts.PATCH(":id", c.UpdateAccount)
accounts.POST(":id/images", c.UploadAccountImage)
```