package models

import (
	"fmt"
	"github.com/kercylan98/kspace/src/pkg/model"
)

// OAuth2Client 基于 OAuth2 的客户端模型
type OAuth2Client struct {
	model.Core

	UserId       uint   `form:"user_id" json:"user_id"`
	ClientID     string `gorm:"uniqueIndex;size:256" form:"client_id" json:"clientID"`
	ClientSecret string `gorm:"size:256" form:"client_secret" json:"clientSecret"`
	Domain       string `gorm:"size:32" form:"domain" json:"domain"`
}

func (slf OAuth2Client) GetID() string {
	return slf.ClientID
}

func (slf OAuth2Client) GetSecret() string {
	return slf.ClientSecret
}

func (slf OAuth2Client) GetDomain() string {
	return slf.Domain
}

func (slf OAuth2Client) GetUserID() string {
	return fmt.Sprint(slf.UserId)
}
