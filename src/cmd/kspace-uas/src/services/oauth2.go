package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	oauth "github.com/go-oauth2/oauth2/v4"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/pkg/models"
	"github.com/kercylan98/kspace/src/cmd/kspace-dal/src/rpc"
	"github.com/kercylan98/kspace/src/cmd/kspace-uas/src/pkg/oauth2"
	"github.com/kercylan98/kspace/src/pkg/web"
	"net/http"
	"time"
)

// OAuth2 参考：https://blog.csdn.net/qq_38384460/article/details/118221000
type OAuth2 struct {
	oauth2.Server[oauth2.Redis]
	rpc.DalOAuth2Client
	rpc.DalUserClient
}

func (slf OAuth2) Initialization(router gin.IRouter, twist web.TwistFunc) {
	oauthGroup := router.Group("/oauth")
	oauthGroup.
		GET("/authorize", twist.Exec(slf.Authorize)).
		POST("/token", twist.Exec(slf.Token)).
		GET("/test", twist.Exec(slf.ValidationBearerToken)).
		GET("/logout", twist.Exec(slf.Logout))
	oauthGroup.Group("/clients").
		POST("", twist.Exec(slf.CreateClient))

	//accounts := v1.Group("/accounts")
	//{
	//	accounts.GET(":id", c.ShowAccount)
	//	accounts.GET("", c.ListAccounts)
	//	accounts.POST("", c.AddAccount)
	//	accounts.DELETE(":id", c.DeleteAccount)
	//	accounts.PATCH(":id", c.UpdateAccount)
	//	accounts.POST(":id/images", c.UploadAccountImage)
	//}
	slf.SetUserAuthorizationHandler(slf.UserAuthorization)
	slf.SetPasswordAuthorizationHandler(slf.PasswordAuthorization)
}

func (slf OAuth2) Logout(ctx web.Context) (response web.Response) {
	token, err := slf.Server.ValidationBearerToken(ctx.Request)
	if err != nil {
		return response.Err(err).Throw()
	}
	if err = slf.Manager.RemoveAccessToken(ctx.Request.Context(), token.GetAccess()); err != nil {
		return response.Err(err).Throw()
	}

	if err = slf.Manager.RemoveRefreshToken(ctx.Request.Context(), token.GetRefresh()); err != nil {
		return response.Err(err).Throw()
	}

	ctx.Redirect(http.StatusPermanentRedirect, ctx.Query("redirect_uri"))
	return response
}

func (slf OAuth2) ValidationBearerToken(ctx web.Context) (response web.Response) {
	token, err := slf.Server.ValidationBearerToken(ctx.Request)
	if err != nil {
		return response.ErrJSON(err).Throw()
	}
	return response.Pass(map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
		"scope":      token.GetScope(),
	})
}

func (slf OAuth2) PasswordAuthorization(ctx context.Context, username, password string) (userID string, err error) {
	reply, err := slf.DalUserClient.VerifyPassword(ctx, &rpc.User{
		Account:  username,
		Password: password,
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return fmt.Sprint(reply.Id), nil
}

func (slf OAuth2) UserAuthorization(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if client, err := slf.Server.Manager.GetClient(r.Context(), r.FormValue("client_id")); err != nil {
		return "", err
	} else {
		return client.GetUserID(), nil
	}
}

func (slf OAuth2) CreateClient(ctx web.Context) (response web.Response) {
	var oa2c = models.OAuth2Client{}
	if err := ctx.ShouldBind(&oa2c); err != nil {
		return response.Err(err).
			MaybeSo("check whether the incoming parameters are as expected").
			Show("客户端创建失败，请稍后再试。").
			Throw()
	}

	reply, err := slf.DalOAuth2Client.CreateClient(ctx.Request.Context(), &rpc.AuthClient{
		UserID:       uint32(oa2c.UserId),
		ClientID:     oa2c.ClientID,
		ClientSecret: oa2c.ClientSecret,
		Domain:       oa2c.Domain,
	})

	if err != nil {
		return response.Err(err).
			Show("客户端创建失败，请稍后再试。").
			Throw()
	}

	(&oa2c).ID = uint(reply.Id)
	return response.Pass(oa2c)
}

func (slf OAuth2) Authorize(ctx web.Context) (response web.Response) {
	if err := slf.HandleAuthorizeRequest(&web.ResponseWriter{
		ResponseWriter: ctx.Writer,
		Response:       &response,
	}, ctx.Request); err != nil {
		return response.Err(err).Throw()
	}

	return response
}

func (slf OAuth2) Token(ctx web.Context) (response web.Response) {
	if err := slf.HandleTokenRequest(&web.ResponseWriter{
		ResponseWriter: ctx.Writer,
		Response:       &response,
	}, ctx.Request); err != nil {
		response.ErrJSON(err)
	}
	return response
}

func (slf OAuth2) GetByID(ctx context.Context, id string) (oauth.ClientInfo, error) {
	client, err := slf.DalOAuth2Client.GetClientWithClientID(ctx, &rpc.AuthClient{
		ClientID: id,
	})
	if err != nil {
		return nil, err
	}
	return models.OAuth2Client{
		UserId:       uint(client.UserID),
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
		Domain:       client.Domain,
	}, nil
}
