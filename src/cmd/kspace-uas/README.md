# KSpace-UAS `用户及认证系统`
对外部客户端提供了基于`HTTP/HTTPS`的用户认证及管理的接口，对内部应用提供了基于`gRPC`的相关功能
****
## OAuth2
基于`go-oauth2`实现了`OAuth2`功能，由于`Client`需要进行数据库存储，所以依赖了项目`space-dal`，而关于`token`的`redis`的读写则由`kspace-uas`自身维护。