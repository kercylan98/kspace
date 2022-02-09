package models

import (
	"fmt"
	"github.com/kercylan98/kspace/playground/internal/err-game-server/server"
)

type User struct {
	Server server.Server
}

func (slf User) AttackBot(bot Bot) {
	fmt.Println("攻击了机器人", bot)
}

func (slf User) Broadcast(msg string) {
	fmt.Println("播报：", msg)
}
