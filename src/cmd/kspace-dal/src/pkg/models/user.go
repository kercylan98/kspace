package models

import "github.com/kercylan98/kspace/src/pkg/model"

// User 用户模型
type User struct {
	model.Core

	Account       string `gorm:"uniqueIndex" form:"account" json:"account"`
	Password      string `gorm:"size:1024" form:"password" json:"-"`
	OAuth2Clients []OAuth2Client

	Token string `gorm:"-" json:"token"`
}
