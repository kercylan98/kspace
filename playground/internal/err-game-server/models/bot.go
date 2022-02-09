package models

import (
	"fmt"
	"github.com/kercylan98/kspace/playground/internal/err-game-server/server"
)

type Bot struct {
	Server server.Server
}

func (slf Bot) Attacked(user User) {
	fmt.Println("被用户攻击", user)
}

func (slf Bot) Death(source User) {
	for _, user := range slf.Server.Online {
		user.Broadcast(fmt.Sprintf("机器人（%v）被用户（%v）杀死了", slf, user))
	}
}
